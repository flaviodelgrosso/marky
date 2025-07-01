package loaders

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"html"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/flaviodelgrosso/marky/internal/mimetypes"
)

type PptxLoader struct{}

// Load reads a PPTX file and converts it to markdown format.
func (*PptxLoader) Load(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read PPTX file: %w", err)
	}

	converter := NewPptxConverter()

	result, err := converter.Convert(data, ConvertOptions{KeepDataURIs: true})
	if err != nil {
		return "", fmt.Errorf("failed to convert PPTX to markdown: %w", err)
	}
	return result.Markdown, nil
}

// CanLoadMimeType returns true if the MIME type is supported for PPTX files.
func (*PptxLoader) CanLoadMimeType(mimeType string) bool {
	supportedTypes := []string{
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.openxmlformats-officedocument.presentationml",
	}
	return mimetypes.IsMimeTypeSupported(mimeType, supportedTypes)
}

// DocumentConverterResult represents the conversion result
type DocumentConverterResult struct {
	Markdown string
}

// PptxConverter handles PPTX to Markdown conversion
type PptxConverter struct {
	acceptedMimeTypePrefixes []string
	acceptedExtensions       []string
}

// ConvertOptions holds configuration for the conversion
type ConvertOptions struct {
	KeepDataURIs bool
}

// NewPptxConverter creates a new PPTX converter
func NewPptxConverter() *PptxConverter {
	return &PptxConverter{
		acceptedMimeTypePrefixes: []string{
			"application/vnd.openxmlformats-officedocument.presentationml",
		},
		acceptedExtensions: []string{".pptx"},
	}
}

// Convert converts PPTX content to Markdown
func (p *PptxConverter) Convert(data []byte, options ConvertOptions) (*DocumentConverterResult, error) {
	reader := bytes.NewReader(data)
	zipReader, err := zip.NewReader(reader, int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("failed to open PPTX file: %w", err)
	}

	presentation, err := p.parsePresentationXML(zipReader)
	if err != nil {
		return nil, fmt.Errorf("failed to parse presentation: %w", err)
	}

	slides := p.parseSlides(zipReader, presentation)

	markdown := p.convertSlidesToMarkdown(slides, zipReader, options)

	return &DocumentConverterResult{
		Markdown: strings.TrimSpace(markdown),
	}, nil
}

// Presentation represents the structure of the PPTX presentation
type Presentation struct {
	SlideIDs []SlideID `xml:"sldIdLst>sldId"`
}

type SlideID struct {
	ID  string `xml:"id,attr"`
	RID string `xml:"r:id,attr"`
}

type Slide struct {
	CommonSlideData CommonSlideData `xml:"cSld"`
	Notes           *Notes          `xml:"notes,omitempty"`
}

type CommonSlideData struct {
	ShapeTree ShapeTree `xml:"spTree"`
}

type ShapeTree struct {
	Shapes []Shape `xml:"sp"`
	Pics   []Pic   `xml:"pic"`
	Tables []Table `xml:"graphicFrame"`
	Groups []Group `xml:"grpSp"`
}

type Shape struct {
	TextBody *TextBody `xml:"txBody"`
	NvSpPr   NvSpPr    `xml:"nvSpPr"`
}

type Pic struct {
	NvPicPr  NvPicPr  `xml:"nvPicPr"`
	BlipFill BlipFill `xml:"blipFill"`
}

type Table struct {
	Graphic Graphic `xml:"graphic"`
}

type Group struct {
	Shapes []Shape `xml:"sp"`
	Pics   []Pic   `xml:"pic"`
	Tables []Table `xml:"graphicFrame"`
}

type TextBody struct {
	Paragraphs []Paragraph `xml:"p"`
}

type Paragraph struct {
	Runs []Run `xml:"r"`
}

type Run struct {
	Text string `xml:"t"`
}

type NvSpPr struct {
	CNvPr CNvPr `xml:"cNvPr"`
}

type NvPicPr struct {
	CNvPr CNvPr `xml:"cNvPr"`
}

type CNvPr struct {
	Name  string `xml:"name,attr"`
	Descr string `xml:"descr,attr"`
}

type BlipFill struct {
	Blip Blip `xml:"blip"`
}

type Blip struct {
	Embed string `xml:"r:embed,attr"`
}

type Graphic struct {
	GraphicData GraphicData `xml:"graphicData"`
}

type GraphicData struct {
	Table TableData `xml:"tbl"`
}

type TableData struct {
	Rows []TableRow `xml:"tr"`
}

type TableRow struct {
	Cells []TableCell `xml:"tc"`
}

type TableCell struct {
	TextBody TextBody `xml:"txBody"`
}

type Notes struct {
	Text string `xml:",innerxml"`
}

func (*PptxConverter) parsePresentationXML(zipReader *zip.Reader) (*Presentation, error) {
	for _, file := range zipReader.File {
		if file.Name == "ppt/presentation.xml" {
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return nil, err
			}

			var presentation Presentation
			err = xml.Unmarshal(data, &presentation)
			if err != nil {
				return nil, err
			}

			return &presentation, nil
		}
	}
	return nil, errors.New("presentation.xml not found")
}

func (p *PptxConverter) parseSlides(zipReader *zip.Reader, presentation *Presentation) []*Slide {
	var slides []*Slide

	for i := range presentation.SlideIDs {
		slideFile := fmt.Sprintf("ppt/slides/slide%d.xml", i+1)

		for _, file := range zipReader.File {
			if file.Name == slideFile {
				rc, err := file.Open()
				if err != nil {
					continue
				}

				data, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					continue
				}

				var slide Slide
				err = xml.Unmarshal(data, &slide)
				if err != nil {
					continue
				}

				// Try to parse notes
				notesFile := fmt.Sprintf("ppt/notesSlides/notesSlide%d.xml", i+1)
				p.parseSlideNotes(zipReader, notesFile, &slide)

				slides = append(slides, &slide)
				break
			}
		}
	}

	return slides
}

func (*PptxConverter) parseSlideNotes(zipReader *zip.Reader, notesFile string, slide *Slide) {
	for _, file := range zipReader.File {
		if file.Name == notesFile {
			rc, err := file.Open()
			if err != nil {
				return
			}
			defer rc.Close()

			data, err := io.ReadAll(rc)
			if err != nil {
				return
			}

			// Simple text extraction from notes
			text := string(data)
			re := regexp.MustCompile(`<a:t>([^<]*)</a:t>`)
			matches := re.FindAllStringSubmatch(text, -1)

			var notesText strings.Builder
			for _, match := range matches {
				if len(match) > 1 {
					notesText.WriteString(match[1])
					notesText.WriteString(" ")
				}
			}

			if notesText.Len() > 0 {
				slide.Notes = &Notes{Text: strings.TrimSpace(notesText.String())}
			}
			break
		}
	}
}

func (p *PptxConverter) convertSlidesToMarkdown(slides []*Slide, zipReader *zip.Reader, options ConvertOptions) string {
	var markdown strings.Builder

	for i, slide := range slides {
		slideNum := i + 1
		markdown.WriteString(fmt.Sprintf("\n\n<!-- Slide number: %d -->\n", slideNum))

		// Process shapes, pictures, and tables
		p.processShapes(slide.CommonSlideData.ShapeTree.Shapes, &markdown, true)
		p.processPics(slide.CommonSlideData.ShapeTree.Pics, &markdown, zipReader, options)
		p.processTables(slide.CommonSlideData.ShapeTree.Tables, &markdown)
		p.processGroups(slide.CommonSlideData.ShapeTree.Groups, &markdown, zipReader, options)

		// Add notes if present
		if slide.Notes != nil && slide.Notes.Text != "" {
			markdown.WriteString("\n\n### Notes:\n")
			markdown.WriteString(slide.Notes.Text)
		}
	}

	return markdown.String()
}

func (p *PptxConverter) processShapes(shapes []Shape, markdown *strings.Builder, isTitle bool) {
	// Sort shapes by position (simplified - just by order for now)
	for _, shape := range shapes {
		if shape.TextBody != nil {
			text := p.extractTextFromTextBody(shape.TextBody)
			if text != "" {
				if isTitle && len(shapes) > 0 {
					markdown.WriteString("# ")
					markdown.WriteString(strings.TrimSpace(text))
					markdown.WriteString("\n")
					isTitle = false // Only first shape with text is title
				} else {
					markdown.WriteString(text)
					markdown.WriteString("\n")
				}
			}
		}
	}
}

func (p *PptxConverter) processPics(pics []Pic, markdown *strings.Builder, zipReader *zip.Reader, options ConvertOptions) {
	for _, pic := range pics {
		altText := pic.NvPicPr.CNvPr.Descr
		if altText == "" {
			altText = pic.NvPicPr.CNvPr.Name
		}

		// Clean alt text
		altText = regexp.MustCompile(`[\r\n\[\]]`).ReplaceAllString(altText, " ")
		altText = regexp.MustCompile(`\s+`).ReplaceAllString(altText, " ")
		altText = strings.TrimSpace(altText)

		if options.KeepDataURIs && pic.BlipFill.Blip.Embed != "" {
			// Try to get the actual image data
			imageData := p.getImageData(zipReader)
			if imageData != nil {
				b64String := base64.StdEncoding.EncodeToString(imageData)
				fmt.Fprintf(markdown, "\n![%s](data:image/png;base64,%s)\n", altText, b64String)
			} else {
				fmt.Fprintf(markdown, "\n![%s](%s.jpg)\n", altText, p.sanitizeFilename(altText))
			}
		} else {
			filename := p.sanitizeFilename(altText) + ".jpg"
			fmt.Fprintf(markdown, "\n![%s](%s)\n", altText, filename)
		}
	}
}

func (p *PptxConverter) processTables(tables []Table, markdown *strings.Builder) {
	for _, table := range tables {
		markdown.WriteString(p.convertTableToMarkdown(table.Graphic.GraphicData.Table))
	}
}

func (p *PptxConverter) processGroups(groups []Group, markdown *strings.Builder, zipReader *zip.Reader, options ConvertOptions) {
	for _, group := range groups {
		p.processShapes(group.Shapes, markdown, false)
		p.processPics(group.Pics, markdown, zipReader, options)
		p.processTables(group.Tables, markdown)
	}
}

func (*PptxConverter) extractTextFromTextBody(textBody *TextBody) string {
	var text strings.Builder

	for _, paragraph := range textBody.Paragraphs {
		for _, run := range paragraph.Runs {
			text.WriteString(run.Text)
		}
		text.WriteString("\n")
	}

	return strings.TrimSpace(text.String())
}

func (p *PptxConverter) convertTableToMarkdown(table TableData) string {
	if len(table.Rows) == 0 {
		return ""
	}

	var markdown strings.Builder

	// Process header row
	if len(table.Rows) > 0 {
		markdown.WriteString("|")
		for _, cell := range table.Rows[0].Cells {
			cellText := p.extractTextFromTextBody(&cell.TextBody)
			cellText = html.EscapeString(cellText)
			markdown.WriteString(" ")
			markdown.WriteString(cellText)
			markdown.WriteString(" |")
		}
		markdown.WriteString("\n")

		// Add separator
		markdown.WriteString("|")
		for range table.Rows[0].Cells {
			markdown.WriteString("---|")
		}
		markdown.WriteString("\n")
	}

	// Process data rows
	for i := 1; i < len(table.Rows); i++ {
		markdown.WriteString("|")
		for _, cell := range table.Rows[i].Cells {
			cellText := p.extractTextFromTextBody(&cell.TextBody)
			cellText = html.EscapeString(cellText)
			markdown.WriteString(" ")
			markdown.WriteString(cellText)
			markdown.WriteString(" |")
		}
		markdown.WriteString("\n")
	}

	return markdown.String()
}

func (*PptxConverter) getImageData(zipReader *zip.Reader) []byte {
	// This is a simplified version - in a real implementation,
	// you'd need to parse the relationship files to map embed IDs to actual files
	for _, file := range zipReader.File {
		if strings.HasPrefix(file.Name, "ppt/media/") {
			rc, err := file.Open()
			if err != nil {
				continue
			}

			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			return data
		}
	}
	return nil
}

func (*PptxConverter) sanitizeFilename(filename string) string {
	re := regexp.MustCompile(`\W`)
	return re.ReplaceAllString(filename, "")
}
