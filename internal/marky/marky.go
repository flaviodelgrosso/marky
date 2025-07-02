package marky

import (
	"fmt"
	"slices"

	"github.com/flaviodelgrosso/marky/internal/converters"
	"github.com/gabriel-vasile/mimetype"
)

// Marky manages document converters and provides conversion functionality.
type Marky struct {
	Converters []converters.Converter
}

type IMarky interface {
	Convert(path string) (string, error)
}

// RegisterConverter adds a new document converter to the available converters.
func (m *Marky) RegisterConverter(converter converters.Converter) {
	m.Converters = append(m.Converters, converter)
}

// Convert processes a document file and converts it to markdown format.
// Returns the markdown content and an error if the conversion fails.
func (m *Marky) Convert(path string) (string, error) {
	// Detect MIME type from file content - this is mandatory
	mtype, err := mimetype.DetectFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to detect MIME type: %w", err)
	}

	// Find a converter that can handle this MIME type
	for _, converter := range m.Converters {
		if accepts(mtype, converter.AcceptedExtensions(), converter.AcceptedMimeTypes()) {
			return converter.Load(path)
		}
	}

	return "", fmt.Errorf("no converter found for MIME type: %s", mtype.String())
}

func accepts(mtype *mimetype.MIME, extensions, mtypes []string) bool {
	if slices.Contains(extensions, mtype.Extension()) {
		return true
	}

	return slices.ContainsFunc(mtypes, mtype.Is)
}
