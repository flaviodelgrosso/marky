package loaders

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/flaviodelgrosso/marky/internal/mimetypes"
	"github.com/flaviodelgrosso/marky/internal/utils"
)

// DocLoader handles loading and converting DOC and DOCX files to markdown.
type DocLoader struct{}

// Load reads a DOC or DOCX file and converts it to markdown.
func (*DocLoader) Load(filePath string) (string, error) {
	content, err := convertDocxToMarkdown(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to convert document: %w", err)
	}

	return content, nil
}

// CanLoadMimeType returns true if the MIME type is supported for DOC/DOCX files.
func (*DocLoader) CanLoadMimeType(mimeType string) bool {
	supportedTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.wordprocessingml",
		"application/msword",
	}
	return mimetypes.IsMimeTypeSupported(mimeType, supportedTypes)
}

// Relationship is
type Relationship struct {
	Text       string `xml:",chardata"`
	ID         string `xml:"Id,attr"`
	Type       string `xml:"Type,attr"`
	Target     string `xml:"Target,attr"`
	TargetMode string `xml:"TargetMode,attr"`
}

// Relationships is
type Relationships struct {
	XMLName      xml.Name       `xml:"Relationships"`
	Text         string         `xml:",chardata"`
	Xmlns        string         `xml:"xmlns,attr"`
	Relationship []Relationship `xml:"Relationship"`
}

// TextVal is
type TextVal struct {
	Text string `xml:",chardata"`
	Val  string `xml:"val,attr"`
}

// NumberingLvl is
type NumberingLvl struct {
	Text      string  `xml:",chardata"`
	Ilvl      string  `xml:"ilvl,attr"`
	Tplc      string  `xml:"tplc,attr"`
	Tentative string  `xml:"tentative,attr"`
	Start     TextVal `xml:"start"`
	NumFmt    TextVal `xml:"numFmt"`
	LvlText   TextVal `xml:"lvlText"`
	LvlJc     TextVal `xml:"lvlJc"`
	PPr       struct {
		Text string `xml:",chardata"`
		Ind  struct {
			Text    string `xml:",chardata"`
			Left    string `xml:"left,attr"`
			Hanging string `xml:"hanging,attr"`
		} `xml:"ind"`
	} `xml:"pPr"`
	RPr struct {
		Text string `xml:",chardata"`
		U    struct {
			Text string `xml:",chardata"`
			Val  string `xml:"val,attr"`
		} `xml:"u"`
		RFonts struct {
			Text string `xml:",chardata"`
			Hint string `xml:"hint,attr"`
		} `xml:"rFonts"`
	} `xml:"rPr"`
}

// Numbering is
type Numbering struct {
	XMLName     xml.Name `xml:"numbering"`
	Text        string   `xml:",chardata"`
	Wpc         string   `xml:"wpc,attr"`
	Cx          string   `xml:"cx,attr"`
	Cx1         string   `xml:"cx1,attr"`
	Mc          string   `xml:"mc,attr"`
	O           string   `xml:"o,attr"`
	R           string   `xml:"r,attr"`
	M           string   `xml:"m,attr"`
	V           string   `xml:"v,attr"`
	Wp14        string   `xml:"wp14,attr"`
	Wp          string   `xml:"wp,attr"`
	W10         string   `xml:"w10,attr"`
	W           string   `xml:"w,attr"`
	W14         string   `xml:"w14,attr"`
	W15         string   `xml:"w15,attr"`
	W16se       string   `xml:"w16se,attr"`
	Wpg         string   `xml:"wpg,attr"`
	Wpi         string   `xml:"wpi,attr"`
	Wne         string   `xml:"wne,attr"`
	Wps         string   `xml:"wps,attr"`
	Ignorable   string   `xml:"Ignorable,attr"`
	AbstractNum []struct {
		Text                       string         `xml:",chardata"`
		AbstractNumID              string         `xml:"abstractNumId,attr"`
		RestartNumberingAfterBreak string         `xml:"restartNumberingAfterBreak,attr"`
		Nsid                       TextVal        `xml:"nsid"`
		MultiLevelType             TextVal        `xml:"multiLevelType"`
		Tmpl                       TextVal        `xml:"tmpl"`
		Lvl                        []NumberingLvl `xml:"lvl"`
	} `xml:"abstractNum"`
	Num []struct {
		Text          string  `xml:",chardata"`
		NumID         string  `xml:"numId,attr"`
		AbstractNumID TextVal `xml:"abstractNumId"`
	} `xml:"num"`
}

type file struct {
	rels  Relationships
	num   Numbering
	r     *zip.ReadCloser
	embed bool
	list  map[string]int
}

// Node is
type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:"-"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

// UnmarshalXML is
func (n *Node) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	n.Attrs = start.Attr
	type node Node

	return d.DecodeElement((*node)(n), &start)
}

func escape(s, set string) string {
	replacer := []string{}
	for _, r := range set {
		rs := string(r)
		replacer = append(replacer, rs, `\`+rs)
	}
	return strings.NewReplacer(replacer...).Replace(s)
}

func (zf *file) extract(rel *Relationship, w io.Writer) error {
	err := os.MkdirAll(filepath.Dir(rel.Target), 0o755)
	if err != nil {
		return err
	}
	for _, f := range zf.r.File {
		if f.Name != "word/"+rel.Target {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		b := make([]byte, f.UncompressedSize64)
		n, err := rc.Read(b)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		if zf.embed {
			fmt.Fprintf(w, "![](data:image/png;base64,%s)",
				base64.StdEncoding.EncodeToString(b[:n]))
		} else {
			err = os.WriteFile(rel.Target, b, 0o644)
			if err != nil {
				return err
			}
			fmt.Fprintf(w, "![](%s)", escape(rel.Target, "()"))
		}
		break
	}
	return nil
}

func attr(attrs []xml.Attr, name string) (string, bool) {
	for _, attr := range attrs {
		if attr.Name.Local == name {
			return attr.Value, true
		}
	}
	return "", false
}

func (zf *file) walk(node *Node, w io.Writer) error {
	switch node.XMLName.Local {
	case "hyperlink":
		return zf.handleHyperlink(node, w)
	case "t":
		fmt.Fprint(w, string(node.Content))
	case "pPr":
		return zf.handlePPr(node, w)
	case "tbl":
		return zf.handleTbl(node, w)
	case "r":
		return zf.handleR(node, w)
	case "p":
		for _, n := range node.Nodes {
			if err := zf.walk(&n, w); err != nil {
				return err
			}
		}
		fmt.Fprintln(w)
	case "blip":
		return zf.handleBlip(node, w)
	case "Fallback":
		// no-op
	case "txbxContent":
		var cbuf bytes.Buffer
		for _, n := range node.Nodes {
			if err := zf.walk(&n, &cbuf); err != nil {
				return err
			}
		}
		fmt.Fprintln(w, "\n```\n"+cbuf.String()+"```")
	default:
		for _, n := range node.Nodes {
			if err := zf.walk(&n, w); err != nil {
				return err
			}
		}
	}
	return nil
}

// --- Helper methods for walk ---

func (zf *file) handleHyperlink(node *Node, w io.Writer) error {
	fmt.Fprint(w, "[")
	var cbuf bytes.Buffer
	for _, n := range node.Nodes {
		if err := zf.walk(&n, &cbuf); err != nil {
			return err
		}
	}
	fmt.Fprint(w, escape(cbuf.String(), "[]"))
	fmt.Fprint(w, "]")

	fmt.Fprint(w, "(")
	if id, ok := attr(node.Attrs, "id"); ok {
		for _, rel := range zf.rels.Relationship {
			if id == rel.ID {
				fmt.Fprint(w, escape(rel.Target, "()"))
				break
			}
		}
	}
	fmt.Fprint(w, ")")
	return nil
}

func (zf *file) handlePPr(node *Node, w io.Writer) error {
	code := zf.processPPrNodes(node, w)

	if code {
		fmt.Fprint(w, "`")
	}
	for _, n := range node.Nodes {
		if err := zf.walk(&n, w); err != nil {
			return err
		}
	}
	if code {
		fmt.Fprint(w, "`")
	}
	return nil
}

func (zf *file) processPPrNodes(node *Node, w io.Writer) bool {
	code := false
	for _, n := range node.Nodes {
		switch n.XMLName.Local {
		case "ind":
			handleIndentation(&n, w)
		case "pStyle":
			if handleParagraphStyle(&n, w) {
				code = true
			}
		case "numPr":
			zf.handleNumPr(&n, w)
		}
	}
	return code
}

func handleIndentation(n *Node, w io.Writer) {
	if left, ok := attr(n.Attrs, "left"); ok {
		if i, err := strconv.Atoi(left); err == nil && i > 0 {
			fmt.Fprint(w, strings.Repeat("  ", i/360))
		}
	}
}

func handleParagraphStyle(n *Node, w io.Writer) bool {
	val, ok := attr(n.Attrs, "val")
	if !ok {
		return false
	}

	switch {
	case strings.HasPrefix(val, "Heading"):
		writeHeading(val, w)
	case val == "Code":
		return true
	default:
		writeNumericHeading(val, w)
	}
	return false
}

func writeHeading(val string, w io.Writer) {
	if i, err := strconv.Atoi(val[7:]); err == nil && i > 0 {
		fmt.Fprint(w, strings.Repeat("#", i)+" ")
	}
}

func writeNumericHeading(val string, w io.Writer) {
	if i, err := strconv.Atoi(val); err == nil && i > 0 {
		fmt.Fprint(w, strings.Repeat("#", i)+" ")
	}
}

func (zf *file) handleNumPr(n *Node, w io.Writer) {
	numID, ilvl := extractNumProperties(n)
	numFmt, start, ind := zf.findNumberingFormat(numID, ilvl)
	zf.writeNumbering(numID, numFmt, start, ind, w)
}

func extractNumProperties(n *Node) (numID, ilvl string) {
	for _, nn := range n.Nodes {
		switch nn.XMLName.Local {
		case "numId":
			if val, ok := attr(nn.Attrs, "val"); ok {
				numID = val
			}
		case "ilvl":
			if val, ok := attr(nn.Attrs, "val"); ok {
				ilvl = val
			}
		}
	}
	return numID, ilvl
}

func (zf *file) findNumberingFormat(numID, ilvl string) (numFmt string, start, ind int) {
	start = 1
	ind = 0

	for _, num := range zf.num.Num {
		if numID != num.NumID {
			continue
		}
		numFmt, start, ind = zf.processAbstractNum(num.AbstractNumID.Val, ilvl)
		break
	}
	return numFmt, start, ind
}

func (zf *file) processAbstractNum(abstractNumID, ilvl string) (numFmt string, start, ind int) {
	start = 1
	ind = 0

	for _, abnum := range zf.num.AbstractNum {
		if abnum.AbstractNumID != abstractNumID {
			continue
		}
		numFmt, start, ind = processAbstractNumLevel(abnum.Lvl, ilvl)
		break
	}
	return numFmt, start, ind
}

func processAbstractNumLevel(levels []NumberingLvl, ilvl string) (numFmt string, start, ind int) {
	start = 1
	ind = 0

	for _, ablvl := range levels {
		if ablvl.Ilvl != ilvl {
			continue
		}
		if i, err := strconv.Atoi(ablvl.Start.Val); err == nil {
			start = i
		}
		if i, err := strconv.Atoi(ablvl.PPr.Ind.Left); err == nil {
			ind = i / 360
		}
		numFmt = ablvl.NumFmt.Val
		break
	}
	return numFmt, start, ind
}

func (zf *file) writeNumbering(numID, numFmt string, start, ind int, w io.Writer) {
	fmt.Fprint(w, strings.Repeat("  ", ind))
	switch numFmt {
	case "decimal", "aiueoFullWidth":
		zf.writeOrderedList(numID, start, ind, w)
	case "bullet":
		fmt.Fprint(w, "* ")
	}
}

func (zf *file) writeOrderedList(numID string, start, ind int, w io.Writer) {
	key := fmt.Sprintf("%s:%d", numID, ind)
	cur, ok := zf.list[key]
	if !ok {
		zf.list[key] = start
	} else {
		zf.list[key] = cur + 1
	}
	fmt.Fprintf(w, "%d. ", zf.list[key])
}

func (zf *file) handleTbl(node *Node, w io.Writer) error {
	rows := zf.extractTableRows(node)
	if len(rows) == 0 {
		return nil
	}

	maxcol := calculateMaxColumns(rows)
	widths := calculateColumnWidths(rows, maxcol)
	writeMarkdownTable(rows, widths, maxcol, w)

	return nil
}

func (zf *file) extractTableRows(node *Node) [][]string {
	var rows [][]string
	for _, tr := range node.Nodes {
		if tr.XMLName.Local != "tr" {
			continue
		}
		cols := zf.extractTableColumns(&tr)
		if len(cols) > 0 {
			rows = append(rows, cols)
		}
	}
	return rows
}

func (zf *file) extractTableColumns(tr *Node) []string {
	// Pre-allocate slice with estimated capacity based on number of child nodes
	cols := make([]string, 0, len(tr.Nodes))
	for _, tc := range tr.Nodes {
		if tc.XMLName.Local != "tc" {
			continue
		}
		var cbuf bytes.Buffer
		if err := zf.walk(&tc, &cbuf); err != nil {
			// Continue processing other columns even if one fails
			cols = append(cols, "")
			continue
		}
		cols = append(cols, strings.ReplaceAll(cbuf.String(), "\n", ""))
	}
	return cols
}

func calculateMaxColumns(rows [][]string) int {
	maxcol := 0
	for _, cols := range rows {
		if len(cols) > maxcol {
			maxcol = len(cols)
		}
	}
	return maxcol
}

func calculateColumnWidths(rows [][]string, maxcol int) []int {
	widths := make([]int, maxcol)
	for _, row := range rows {
		for i := range maxcol {
			if i < len(row) {
				width := utils.StringWidth(row[i])
				if widths[i] < width {
					widths[i] = width
				}
			}
		}
	}
	return widths
}

func writeMarkdownTable(rows [][]string, widths []int, maxcol int, w io.Writer) {
	for i, row := range rows {
		if i == 0 {
			writeTableHeader(widths, maxcol, w)
		}
		writeTableRow(row, widths, maxcol, w)
	}
	fmt.Fprint(w, "\n")
}

func writeTableHeader(widths []int, maxcol int, w io.Writer) {
	// Write empty header row
	for j := 0; j < maxcol; j++ {
		fmt.Fprint(w, "|")
		fmt.Fprint(w, strings.Repeat(" ", widths[j]))
	}
	fmt.Fprint(w, "|\n")

	// Write separator row
	for j := range maxcol {
		fmt.Fprint(w, "|")
		fmt.Fprint(w, strings.Repeat("-", widths[j]))
	}
	fmt.Fprint(w, "|\n")
}

func writeTableRow(row []string, widths []int, maxcol int, w io.Writer) {
	for j := range maxcol {
		fmt.Fprint(w, "|")
		if j < len(row) {
			width := utils.StringWidth(row[j])
			fmt.Fprint(w, escape(row[j], "|"))
			fmt.Fprint(w, strings.Repeat(" ", widths[j]-width))
		} else {
			fmt.Fprint(w, strings.Repeat(" ", widths[j]))
		}
	}
	fmt.Fprint(w, "|\n")
}

func (zf *file) handleR(node *Node, w io.Writer) error {
	bold := false
	italic := false
	strike := false
	for _, n := range node.Nodes {
		if n.XMLName.Local != "rPr" {
			continue
		}
		for _, nn := range n.Nodes {
			switch nn.XMLName.Local {
			case "b":
				bold = true
			case "i":
				italic = true
			case "strike":
				strike = true
			}
		}
	}
	if strike {
		fmt.Fprint(w, "~~")
	}
	if bold {
		fmt.Fprint(w, "**")
	}
	if italic {
		fmt.Fprint(w, "*")
	}
	var cbuf bytes.Buffer
	for _, n := range node.Nodes {
		if err := zf.walk(&n, &cbuf); err != nil {
			return err
		}
	}
	fmt.Fprint(w, escape(cbuf.String(), `*~\`))
	if italic {
		fmt.Fprint(w, "*")
	}
	if bold {
		fmt.Fprint(w, "**")
	}
	if strike {
		fmt.Fprint(w, "~~")
	}
	return nil
}

func (zf *file) handleBlip(node *Node, w io.Writer) error {
	if id, ok := attr(node.Attrs, "embed"); ok {
		for _, rel := range zf.rels.Relationship {
			if id != rel.ID {
				continue
			}
			if err := zf.extract(&rel, w); err != nil {
				return err
			}
		}
	}
	return nil
}

func readDocFile(f *zip.File) (*Node, error) {
	rc, err := f.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		return nil, err
	}

	var node Node

	err = xml.Unmarshal(b, &node)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func findFile(files []*zip.File, target string) *zip.File {
	for _, f := range files {
		if ok, _ := path.Match(target, f.Name); ok {
			return f
		}
	}
	return nil
}

func convertDocxToMarkdown(filePath string) (string, error) {
	r, err := zip.OpenReader(filePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var rels Relationships
	var num Numbering

	for _, f := range r.File {
		switch f.Name {
		case "word/_rels/document.xml.rels", "word/_rels/document2.xml.rels":
			rc, err := f.Open()
			defer rc.Close()

			b, _ := io.ReadAll(rc)
			if err != nil {
				return "", err
			}

			err = xml.Unmarshal(b, &rels)
			if err != nil {
				return "", err
			}
		case "word/numbering.xml":
			rc, err := f.Open()
			defer rc.Close()

			b, _ := io.ReadAll(rc)
			if err != nil {
				return "", err
			}

			err = xml.Unmarshal(b, &num)
			if err != nil {
				return "", err
			}
		}
	}

	f := findFile(r.File, "word/document*.xml")
	if f == nil {
		return "", errors.New("incorrect document")
	}
	node, err := readDocFile(f)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	zf := &file{
		r:    r,
		rels: rels,
		num:  num,
		list: make(map[string]int),
	}
	err = zf.walk(node, &buf)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
