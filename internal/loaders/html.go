package loaders

import (
	"fmt"
	"os"

	html2md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/flaviodelgrosso/marky/internal/mimetypes"
)

// HTMLLoader handles loading and converting HTML files to markdown.
type HTMLLoader struct{}

// Load reads an HTML file and converts it to markdown.
func (*HTMLLoader) Load(path string) (string, error) {
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

// CanLoadMimeType returns true if the MIME type is supported for HTML files.
func (*HTMLLoader) CanLoadMimeType(mimeType string) bool {
	supportedTypes := []string{
		"text/html",
	}
	return mimetypes.IsMimeTypeSupported(mimeType, supportedTypes)
}
