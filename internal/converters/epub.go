package converters

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// EpubConverter handles loading and converting EPUB files to markdown.
type EpubConverter struct {
	BaseConverter
}

// NewEpubConverter creates a new EPUB converter with appropriate MIME types and extensions.
func NewEpubConverter() Converter {
	return &EpubConverter{
		BaseConverter: NewBaseConverter(
			[]string{".epub"},
			[]string{
				"application/epub",
				"application/epub+zip",
				"application/x-epub+zip",
			},
		),
	}
}

// Container represents the META-INF/container.xml structure
type Container struct {
	Rootfiles []Rootfile `xml:"rootfiles>rootfile"`
}

type Rootfile struct {
	FullPath  string `xml:"full-path,attr"`
	MediaType string `xml:"media-type,attr"`
}

// Package represents the OPF package structure
type Package struct {
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Metadata struct {
	Title       []string `xml:"title"`
	Creator     []string `xml:"creator"`
	Language    string   `xml:"language"`
	Publisher   string   `xml:"publisher"`
	Date        string   `xml:"date"`
	Description string   `xml:"description"`
	Identifier  string   `xml:"identifier"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	Items []SpineItem `xml:"itemref"`
}

type SpineItem struct {
	IDRef string `xml:"idref,attr"`
}

// Load reads an EPUB file and converts it to markdown.
func (*EpubConverter) Load(path string) (string, error) {
	// Open the EPUB file as a ZIP archive
	reader, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("failed to open EPUB file: %w", err)
	}
	defer reader.Close()

	// Find and parse container.xml
	containerFile, err := findFileInZip(&reader.Reader, "META-INF/container.xml")
	if err != nil {
		return "", fmt.Errorf("failed to find container.xml: %w", err)
	}

	var container Container
	if err := parseXMLFile(containerFile, &container); err != nil {
		return "", fmt.Errorf("failed to parse container.xml: %w", err)
	}

	if len(container.Rootfiles) == 0 {
		return "", errors.New("no rootfiles found in container.xml")
	}

	// Parse the OPF file
	opfPath := container.Rootfiles[0].FullPath
	opfFile, err := findFileInZip(&reader.Reader, opfPath)
	if err != nil {
		return "", fmt.Errorf("failed to find OPF file %s: %w", opfPath, err)
	}

	var pkg Package
	if err := parseXMLFile(opfFile, &pkg); err != nil {
		return "", fmt.Errorf("failed to parse OPF file: %w", err)
	}

	// Create a map of item IDs to hrefs
	manifestMap := make(map[string]string)
	for _, item := range pkg.Manifest.Items {
		manifestMap[item.ID] = item.Href
	}

	// Get the base directory of the OPF file
	baseDir := filepath.Dir(opfPath)

	// Process spine items in order
	var markdownParts []string

	// Add metadata as header
	metadata := formatMetadata(pkg.Metadata)
	if metadata != "" {
		markdownParts = append(markdownParts, metadata)
	}

	// Convert content files
	for _, spineItem := range pkg.Spine.Items {
		href, exists := manifestMap[spineItem.IDRef]
		if !exists {
			continue
		}

		if baseDir != "." && baseDir != "" {
			href = filepath.Join(baseDir, href)
		}

		// Find and convert the content file
		contentFile, err := findFileInZip(&reader.Reader, href)
		if err != nil {
			// Skip missing files
			continue
		}

		markdown, err := convertHTMLToMarkdown(contentFile)
		if err != nil {
			// Skip files that can't be converted
			continue
		}

		if strings.TrimSpace(markdown) != "" {
			markdownParts = append(markdownParts, markdown)
		}
	}

	return strings.Join(markdownParts, "\n\n"), nil
}

func findFileInZip(reader *zip.Reader, filename string) (*zip.File, error) {
	for _, file := range reader.File {
		if file.Name == filename {
			return file, nil
		}
	}
	return nil, fmt.Errorf("file %s not found in ZIP archive", filename)
}

func parseXMLFile(file *zip.File, v any) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return err
	}

	return xml.Unmarshal(data, v)
}

func convertHTMLToMarkdown(file *zip.File) (string, error) {
	rc, err := file.Open()
	if err != nil {
		return "", err
	}
	defer rc.Close()

	content, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}

	// Convert HTML to Markdown
	markdown, err := html2md.ConvertString(string(content))
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(markdown), nil
}

func formatMetadata(metadata Metadata) string {
	var parts []string

	if len(metadata.Title) > 0 && metadata.Title[0] != "" {
		parts = append(parts, "**Title:** "+metadata.Title[0])
	}

	if len(metadata.Creator) > 0 {
		// Filter out empty creators
		var creators []string
		for _, creator := range metadata.Creator {
			if creator != "" {
				creators = append(creators, creator)
			}
		}
		if len(creators) > 0 {
			parts = append(parts, "**Authors:** "+strings.Join(creators, ", "))
		}
	}

	if metadata.Language != "" {
		parts = append(parts, "**Language:** "+metadata.Language)
	}

	if metadata.Publisher != "" {
		parts = append(parts, "**Publisher:** "+metadata.Publisher)
	}

	if metadata.Date != "" {
		parts = append(parts, "**Date:** "+metadata.Date)
	}

	if metadata.Description != "" {
		parts = append(parts, "**Description:** "+metadata.Description)
	}

	if metadata.Identifier != "" {
		parts = append(parts, "**Identifier:** "+metadata.Identifier)
	}

	if len(parts) == 0 {
		return ""
	}

	return strings.Join(parts, "\n")
}
