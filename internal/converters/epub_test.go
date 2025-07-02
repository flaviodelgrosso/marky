package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNewEpubConverter(t *testing.T) {
	converter := NewEpubConverter()

	expectedExtensions := []string{".epub"}
	expectedMimeTypes := []string{
		"application/epub",
		"application/epub+zip",
		"application/x-epub+zip",
	}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewEpubConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewEpubConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestEpubConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewEpubConverter()
	_, err := converter.Load("/nonexistent/file.epub")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	// The error message will depend on the internal implementation
	if err.Error() == "" {
		t.Errorf("Load() should return a meaningful error message")
	}
}

func TestEpubConverter_Load_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid EPUB
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.epub")

	err := os.WriteFile(invalidFile, []byte("not an epub file"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewEpubConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid EPUB file")
	}
}

func TestEpubConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.epub")

	err := os.WriteFile(emptyFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	converter := NewEpubConverter()
	_, err = converter.Load(emptyFile)

	if err == nil {
		t.Errorf("Load() should return error for empty file")
	}
}

func TestEpubConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	noPermFile := filepath.Join(tempDir, "noperm.epub")

	err := os.WriteFile(noPermFile, []byte("fake epub content"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	converter := NewEpubConverter()
	_, err = converter.Load(noPermFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestEpubConverter_Interface(t *testing.T) {
	converter := NewEpubConverter()

	// Test AcceptedExtensions
	extensions := converter.AcceptedExtensions()
	expectedExtensions := []string{".epub"}
	if !reflect.DeepEqual(extensions, expectedExtensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", extensions, expectedExtensions)
	}

	// Test AcceptedMimeTypes
	mimeTypes := converter.AcceptedMimeTypes()
	if len(mimeTypes) != 3 {
		t.Errorf("AcceptedMimeTypes() should return 3 MIME types, got %d", len(mimeTypes))
	}
}

// Test the EPUB-specific structs
func TestContainerStruct(t *testing.T) {
	container := Container{
		Rootfiles: []Rootfile{
			{
				FullPath:  "OEBPS/content.opf",
				MediaType: "application/oebps-package+xml",
			},
		},
	}

	if len(container.Rootfiles) != 1 {
		t.Errorf("Container should have 1 rootfile, got %d", len(container.Rootfiles))
	}

	if container.Rootfiles[0].FullPath != "OEBPS/content.opf" {
		t.Errorf("Rootfile.FullPath = %v, want 'OEBPS/content.opf'", container.Rootfiles[0].FullPath)
	}

	if container.Rootfiles[0].MediaType != "application/oebps-package+xml" {
		t.Errorf("Rootfile.MediaType = %v, want 'application/oebps-package+xml'", container.Rootfiles[0].MediaType)
	}
}

func TestRootfileStruct(t *testing.T) {
	rootfile := Rootfile{
		FullPath:  "content/book.opf",
		MediaType: "application/oebps-package+xml",
	}

	if rootfile.FullPath != "content/book.opf" {
		t.Errorf("Rootfile.FullPath = %v, want 'content/book.opf'", rootfile.FullPath)
	}

	if rootfile.MediaType != "application/oebps-package+xml" {
		t.Errorf("Rootfile.MediaType = %v, want 'application/oebps-package+xml'", rootfile.MediaType)
	}
}

func TestPackageStruct(t *testing.T) {
	// Test basic structure of Package
	pkg := Package{
		Metadata: Metadata{},
		Manifest: Manifest{},
		Spine:    Spine{},
	}

	// Just verify the struct can be created without errors
	// We can't compare structs with slices directly, so we just verify creation works
	_ = pkg // Use the variable to avoid unused variable error
}

// MockEpubConverter for testing scenarios where we need to simulate EPUB processing
type MockEpubConverter struct {
	BaseConverter
	loadFunc func(path string) (string, error)
}

func (m *MockEpubConverter) Load(path string) (string, error) {
	if m.loadFunc != nil {
		return m.loadFunc(path)
	}
	return "Mock EPUB content extracted", nil
}

func TestMockEpubConverter_Success(t *testing.T) {
	mockConverter := &MockEpubConverter{
		BaseConverter: NewBaseConverter([]string{".epub"},
			[]string{"application/epub+zip"}),
		loadFunc: func(path string) (string, error) {
			return "# Book Title\n\nThis is extracted text from EPUB: " + path, nil
		},
	}

	result, err := mockConverter.Load("test.epub")
	if err != nil {
		t.Errorf("Mock Load() returned unexpected error: %v", err)
	}

	expected := "# Book Title\n\nThis is extracted text from EPUB: test.epub"
	if result != expected {
		t.Errorf("Mock Load() = %v, want %v", result, expected)
	}
}

func TestMockEpubConverter_Error(t *testing.T) {
	mockConverter := &MockEpubConverter{
		BaseConverter: NewBaseConverter([]string{".epub"},
			[]string{"application/epub+zip"}),
		loadFunc: func(path string) (string, error) {
			return "", os.ErrNotExist
		},
	}

	_, err := mockConverter.Load("nonexistent.epub")
	if err == nil {
		t.Errorf("Mock Load() should return error")
	}
}

// Test various error scenarios
func TestEpubConverter_LoadingScenarios(t *testing.T) {
	converter := NewEpubConverter()

	testCases := []struct {
		name        string
		filename    string
		shouldError bool
	}{
		{
			name:        "Non-existent file",
			filename:    "/tmp/nonexistent-file-12345.epub",
			shouldError: true,
		},
		{
			name:        "Empty filename",
			filename:    "",
			shouldError: true,
		},
		{
			name:        "Directory instead of file",
			filename:    "/tmp",
			shouldError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := converter.Load(tc.filename)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Load() should return error for %s", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Load() should not return error for %s, got: %v", tc.name, err)
				}
			}
		})
	}
}
