package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewPdfConverter(t *testing.T) {
	converter := NewPdfConverter()

	expectedExtensions := []string{".pdf"}
	expectedMimeTypes := []string{"application/pdf"}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewPdfConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewPdfConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestPdfConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewPdfConverter()
	_, err := converter.Load("/nonexistent/file.pdf")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "unable to open PDF file") {
		t.Errorf("Load() error should mention PDF file opening failure")
	}
}

func TestPdfConverter_Load_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid PDF
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.pdf")

	err := os.WriteFile(invalidFile, []byte("not a pdf file"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewPdfConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid PDF file")
	}

	if !strings.Contains(err.Error(), "unable to open PDF file") {
		t.Errorf("Load() error should mention PDF file opening failure")
	}
}

func TestPdfConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.pdf")

	err := os.WriteFile(emptyFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	converter := NewPdfConverter()
	_, err = converter.Load(emptyFile)

	if err == nil {
		t.Errorf("Load() should return error for empty file")
	}

	if !strings.Contains(err.Error(), "unable to open PDF file") {
		t.Errorf("Load() error should mention PDF file opening failure")
	}
}

func TestPdfConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	noPermFile := filepath.Join(tempDir, "noperm.pdf")

	err := os.WriteFile(noPermFile, []byte("fake pdf content"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	converter := NewPdfConverter()
	_, err = converter.Load(noPermFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestReadPdfFile_NonExistentFile(t *testing.T) {
	_, err := readPdfFile("/nonexistent/file.pdf")

	if err == nil {
		t.Errorf("readPdfFile() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "unable to open PDF file") {
		t.Errorf("readPdfFile() error should mention PDF file opening failure")
	}
}

func TestReadPdfFile_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid PDF
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.pdf")

	err := os.WriteFile(invalidFile, []byte("not a valid pdf"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	_, err = readPdfFile(invalidFile)

	if err == nil {
		t.Errorf("readPdfFile() should return error for invalid PDF file")
	}

	if !strings.Contains(err.Error(), "unable to open PDF file") {
		t.Errorf("readPdfFile() error should mention PDF file opening failure")
	}
}

// Note: Testing with actual PDF files is challenging without creating binary PDF content
// In a real-world scenario, you would either:
// 1. Include test PDF files in your repository
// 2. Generate minimal PDF files programmatically using a PDF library
// 3. Use integration tests with known PDF files
//
// For now, we focus on testing error conditions and the interface compliance.

func TestPdfConverter_Interface(t *testing.T) {
	converter := NewPdfConverter()

	// Test AcceptedExtensions
	extensions := converter.AcceptedExtensions()
	if len(extensions) != 1 || extensions[0] != ".pdf" {
		t.Errorf("AcceptedExtensions() = %v, want [.pdf]", extensions)
	}

	// Test AcceptedMimeTypes
	mimeTypes := converter.AcceptedMimeTypes()
	if len(mimeTypes) != 1 || mimeTypes[0] != "application/pdf" {
		t.Errorf("AcceptedMimeTypes() = %v, want [application/pdf]", mimeTypes)
	}
}

// MockPdfConverter for testing scenarios where we need to simulate PDF processing
type MockPdfConverter struct {
	BaseConverter
	loadFunc func(path string) (string, error)
}

func (m *MockPdfConverter) Load(path string) (string, error) {
	if m.loadFunc != nil {
		return m.loadFunc(path)
	}
	return "Mock PDF content extracted", nil
}

func TestMockPdfConverter_Success(t *testing.T) {
	mockConverter := &MockPdfConverter{
		BaseConverter: NewBaseConverter([]string{".pdf"}, []string{"application/pdf"}),
		loadFunc: func(path string) (string, error) {
			return "This is extracted text from PDF: " + path, nil
		},
	}

	result, err := mockConverter.Load("test.pdf")
	if err != nil {
		t.Errorf("Mock Load() returned unexpected error: %v", err)
	}

	expected := "This is extracted text from PDF: test.pdf"
	if result != expected {
		t.Errorf("Mock Load() = %v, want %v", result, expected)
	}
}

func TestMockPdfConverter_Error(t *testing.T) {
	mockConverter := &MockPdfConverter{
		BaseConverter: NewBaseConverter([]string{".pdf"}, []string{"application/pdf"}),
		loadFunc: func(path string) (string, error) {
			return "", os.ErrNotExist
		},
	}

	_, err := mockConverter.Load("nonexistent.pdf")
	if err == nil {
		t.Errorf("Mock Load() should return error")
	}
}

// Test helper functions for PDF processing scenarios
func TestPdfConverter_LoadingScenarios(t *testing.T) {
	converter := NewPdfConverter()

	testCases := []struct {
		name          string
		filename      string
		shouldError   bool
		errorContains string
	}{
		{
			name:          "Non-existent file",
			filename:      "/tmp/nonexistent-file-12345.pdf",
			shouldError:   true,
			errorContains: "unable to open PDF file",
		},
		{
			name:          "Empty filename",
			filename:      "",
			shouldError:   true,
			errorContains: "unable to open PDF file",
		},
		{
			name:          "Directory instead of file",
			filename:      "/tmp",
			shouldError:   true,
			errorContains: "unable to open PDF file",
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
