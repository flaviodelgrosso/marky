package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewPptxConverter(t *testing.T) {
	converter := NewPptxConverter()

	expectedExtensions := []string{".pptx"}
	expectedMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		"application/vnd.openxmlformats-officedocument.presentationml",
	}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewPptxConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewPptxConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestPptxConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewPptxConverter()
	_, err := converter.Load("/nonexistent/file.pptx")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read PPTX file") {
		t.Errorf("Load() error should mention PPTX file reading failure")
	}
}

func TestPptxConverter_Load_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid PPTX
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.pptx")

	err := os.WriteFile(invalidFile, []byte("not a pptx file"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewPptxConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid PPTX file")
	}

	if !strings.Contains(err.Error(), "failed to convert PPTX to markdown") {
		t.Errorf("Load() error should mention PPTX conversion failure")
	}
}

func TestPptxConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.pptx")

	err := os.WriteFile(emptyFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	converter := NewPptxConverter()
	_, err = converter.Load(emptyFile)

	if err == nil {
		t.Errorf("Load() should return error for empty file")
	}

	if !strings.Contains(err.Error(), "failed to convert PPTX to markdown") {
		t.Errorf("Load() error should mention PPTX conversion failure")
	}
}

func TestPptxConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	noPermFile := filepath.Join(tempDir, "noperm.pptx")

	err := os.WriteFile(noPermFile, []byte("fake pptx content"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	converter := NewPptxConverter()
	_, err = converter.Load(noPermFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestPptxConverter_Interface(t *testing.T) {
	converter := NewPptxConverter()

	// Test AcceptedExtensions
	extensions := converter.AcceptedExtensions()
	expectedExtensions := []string{".pptx"}
	if !reflect.DeepEqual(extensions, expectedExtensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", extensions, expectedExtensions)
	}

	// Test AcceptedMimeTypes
	mimeTypes := converter.AcceptedMimeTypes()
	if len(mimeTypes) != 2 {
		t.Errorf("AcceptedMimeTypes() should return 2 MIME types, got %d", len(mimeTypes))
	}
}

// MockPptxConverter for testing scenarios where we need to simulate PPTX processing
type MockPptxConverter struct {
	BaseConverter
	loadFunc func(path string) (string, error)
}

func (m *MockPptxConverter) Load(path string) (string, error) {
	if m.loadFunc != nil {
		return m.loadFunc(path)
	}
	return "Mock PPTX content extracted", nil
}

func TestMockPptxConverter_Success(t *testing.T) {
	mockConverter := &MockPptxConverter{
		BaseConverter: NewBaseConverter([]string{".pptx"},
			[]string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"}),
		loadFunc: func(path string) (string, error) {
			return "# Slide 1\n\nThis is extracted text from PPTX: " + path, nil
		},
	}

	result, err := mockConverter.Load("test.pptx")
	if err != nil {
		t.Errorf("Mock Load() returned unexpected error: %v", err)
	}

	expected := "# Slide 1\n\nThis is extracted text from PPTX: test.pptx"
	if result != expected {
		t.Errorf("Mock Load() = %v, want %v", result, expected)
	}
}

func TestMockPptxConverter_Error(t *testing.T) {
	mockConverter := &MockPptxConverter{
		BaseConverter: NewBaseConverter([]string{".pptx"},
			[]string{"application/vnd.openxmlformats-officedocument.presentationml.presentation"}),
		loadFunc: func(path string) (string, error) {
			return "", os.ErrNotExist
		},
	}

	_, err := mockConverter.Load("nonexistent.pptx")
	if err == nil {
		t.Errorf("Mock Load() should return error")
	}
}

// Test various error scenarios
func TestPptxConverter_LoadingScenarios(t *testing.T) {
	converter := NewPptxConverter()

	testCases := []struct {
		name          string
		filename      string
		shouldError   bool
		errorContains string
	}{
		{
			name:          "Non-existent file",
			filename:      "/tmp/nonexistent-file-12345.pptx",
			shouldError:   true,
			errorContains: "failed to read PPTX file",
		},
		{
			name:          "Empty filename",
			filename:      "",
			shouldError:   true,
			errorContains: "failed to read PPTX file",
		},
		{
			name:          "Directory instead of file",
			filename:      "/tmp",
			shouldError:   true,
			errorContains: "failed to read PPTX file",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := converter.Load(tc.filename)

			if tc.shouldError {
				if err == nil {
					t.Errorf("Load() should return error for %s", tc.name)
				} else if !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Load() error should contain '%s', got: %v", tc.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Load() should not return error for %s, got: %v", tc.name, err)
				}
			}
		})
	}
}

// Test the DocumentConverterResult struct
func TestDocumentConverterResult(t *testing.T) {
	result := DocumentConverterResult{
		Markdown: "# Test Slide\n\nContent here",
	}

	if result.Markdown != "# Test Slide\n\nContent here" {
		t.Errorf("DocumentConverterResult.Markdown = %v, want '# Test Slide\\n\\nContent here'", result.Markdown)
	}
}

// Test ConvertOptions struct
func TestConvertOptions(t *testing.T) {
	options := ConvertOptions{
		KeepDataURIs: true,
	}

	if !options.KeepDataURIs {
		t.Errorf("ConvertOptions.KeepDataURIs = %v, want true", options.KeepDataURIs)
	}

	options2 := ConvertOptions{
		KeepDataURIs: false,
	}

	if options2.KeepDataURIs {
		t.Errorf("ConvertOptions.KeepDataURIs = %v, want false", options2.KeepDataURIs)
	}
}
