package converters

import (
	"fmt"
	"os"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
)

// HTMLConverter handles loading and converting HTML files to markdown.
type HTMLConverter struct {
	BaseConverter
}

// NewHTMLConverter creates a new HTML converter with appropriate MIME types and extensions.
func NewHTMLConverter() Converter {
	return &HTMLConverter{
		BaseConverter: NewBaseConverter(
			[]string{".html", ".htm"},
			[]string{"text/html"},
		),
	}
}

// Load reads an HTML file and converts it to markdown.
func (*HTMLConverter) Load(path string) (string, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read HTML file: %w", err)
	}

	markdown, err := html2md.ConvertString(string(input))
	if err != nil {
		return "", fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	return markdown, nil
}
