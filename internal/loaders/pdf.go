package loaders

import (
	"bytes"
	"fmt"

	"github.com/flaviodelgrosso/marky/internal/mimetypes"
	"github.com/ledongthuc/pdf"
)

// PdfLoader handles loading and converting PDF files to text.
type PdfLoader struct{}

// Load reads a PDF file and extracts its text content.
func (*PdfLoader) Load(path string) (string, error) {
	return readPdfFile(path)
}

// CanLoadMimeType returns true if the MIME type is supported for PDF files.
func (*PdfLoader) CanLoadMimeType(mimeType string) bool {
	supportedTypes := []string{
		"application/pdf",
	}
	return mimetypes.IsMimeTypeSupported(mimeType, supportedTypes)
}

// readPdfFile reads and extracts text content from a PDF file.
func readPdfFile(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("unable to open PDF file %s: %w", path, err)
	}
	defer f.Close()

	b, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("unable to extract text from PDF file %s: %w", path, err)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(b)
	if err != nil {
		return "", fmt.Errorf("unable to read text content from PDF file %s: %w", path, err)
	}

	return buf.String(), nil
}
