package marky

import (
	"fmt"

	"github.com/flaviodelgrosso/marky/internal/loaders"
	"github.com/flaviodelgrosso/marky/internal/mimetypes"
)

// Marky manages document loaders and provides conversion functionality.
type Marky struct {
	Loaders []loaders.DocumentLoader
}

type IMarky interface {
	Convert(path string) (string, error)
}

// RegisterLoader adds a new document loader to the available loaders.
func (m *Marky) RegisterLoader(loader loaders.DocumentLoader) {
	m.Loaders = append(m.Loaders, loader)
}

// Convert processes a document file and converts it to markdown format.
// Returns the markdown content and an error if the conversion fails.
func (m *Marky) Convert(path string) (string, error) {
	// Detect MIME type from file content - this is mandatory
	mimeInfo, err := mimetypes.DetectMimeType(path)
	if err != nil {
		return "", fmt.Errorf("failed to detect MIME type: %w", err)
	}

	// Find a loader that can handle this MIME type
	for _, loader := range m.Loaders {
		if loader.CanLoadMimeType(mimeInfo.MimeType) {
			return loader.Load(path)
		}
	}

	return "", fmt.Errorf("no loader found for MIME type %s", mimeInfo.MimeType)
}
