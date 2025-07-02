package converters

import (
	"bytes"
	"fmt"

	"github.com/ledongthuc/pdf"
)

// PdfConverter handles loading and converting PDF files to text.
type PdfConverter struct {
	BaseConverter
}

// NewPdfConverter creates a new PDF converter with appropriate MIME types and extensions.
func NewPdfConverter() Converter {
	return &PdfConverter{
		BaseConverter: NewBaseConverter(
			[]string{".pdf"},
			[]string{"application/pdf"},
		),
	}
}

// Load reads a PDF file and extracts its text content.
func (*PdfConverter) Load(path string) (string, error) {
	return readPdfFile(path)
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
