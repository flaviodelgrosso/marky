package loaders

// DocumentLoader defines the interface for loading and converting documents to markdown.
type DocumentLoader interface {
	// Load converts a document at the given path to markdown format.
	// Returns the markdown content and an error if the operation fails.
	Load(path string) (string, error)

	// CanLoadMimeType returns true if this loader can handle files with the given MIME type.
	CanLoadMimeType(mimeType string) bool
}
