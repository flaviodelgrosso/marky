package converters

// Converter defines the interface for document converters.
// It combines metadata about accepted formats with the conversion capability.
type Converter interface {
	// AcceptedMimeTypes returns the MIME types this converter can handle.
	AcceptedMimeTypes() []string

	// AcceptedExtensions returns the file extensions this converter can handle.
	AcceptedExtensions() []string

	// Load converts a document at the given path to markdown format.
	// Returns the markdown content and an error if the operation fails.
	Load(path string) (string, error)
}

// BaseConverter provides a foundation for implementing converters.
// It handles the common metadata fields that all converters need.
type BaseConverter struct {
	acceptedExtensions []string
	acceptedMimeTypes  []string
}

// NewBaseConverter creates a new BaseConverter with the specified extensions and MIME types.
func NewBaseConverter(extensions []string, mimeTypes []string) BaseConverter {
	return BaseConverter{
		acceptedExtensions: extensions,
		acceptedMimeTypes:  mimeTypes,
	}
}

// AcceptedExtensions returns the file extensions this converter can handle.
func (b BaseConverter) AcceptedExtensions() []string {
	return b.acceptedExtensions
}

// AcceptedMimeTypes returns the MIME types this converter can handle.
func (b BaseConverter) AcceptedMimeTypes() []string {
	return b.acceptedMimeTypes
}
