package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewDocConverter(t *testing.T) {
	converter := NewDocConverter()

	expectedExtensions := []string{".docx", ".doc"}
	expectedMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/vnd.openxmlformats-officedocument.wordprocessingml",
		"application/msword",
	}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewDocConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewDocConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestDocConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewDocConverter()
	_, err := converter.Load("/nonexistent/file.docx")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to convert document") {
		t.Errorf("Load() error should mention document conversion failure")
	}
}

func TestDocConverter_Load_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid DOCX
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.docx")

	err := os.WriteFile(invalidFile, []byte("not a docx file"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewDocConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid DOCX file")
	}

	if !strings.Contains(err.Error(), "failed to convert document") {
		t.Errorf("Load() error should mention document conversion failure")
	}
}

func TestDocConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty file
	tempDir := t.TempDir()
	emptyFile := filepath.Join(tempDir, "empty.docx")

	err := os.WriteFile(emptyFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	converter := NewDocConverter()
	_, err = converter.Load(emptyFile)

	if err == nil {
		t.Errorf("Load() should return error for empty file")
	}

	if !strings.Contains(err.Error(), "failed to convert document") {
		t.Errorf("Load() error should mention document conversion failure")
	}
}

func TestDocConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	noPermFile := filepath.Join(tempDir, "noperm.docx")

	err := os.WriteFile(noPermFile, []byte("fake docx content"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	converter := NewDocConverter()
	_, err = converter.Load(noPermFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestDocConverter_Interface(t *testing.T) {
	converter := NewDocConverter()

	// Test AcceptedExtensions
	extensions := converter.AcceptedExtensions()
	expectedExtensions := []string{".docx", ".doc"}
	if !reflect.DeepEqual(extensions, expectedExtensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", extensions, expectedExtensions)
	}

	// Test AcceptedMimeTypes
	mimeTypes := converter.AcceptedMimeTypes()
	if len(mimeTypes) != 3 {
		t.Errorf("AcceptedMimeTypes() should return 3 MIME types, got %d", len(mimeTypes))
	}
}

// MockDocConverter for testing scenarios where we need to simulate DOCX processing
type MockDocConverter struct {
	BaseConverter
	loadFunc func(path string) (string, error)
}

func (m *MockDocConverter) Load(path string) (string, error) {
	if m.loadFunc != nil {
		return m.loadFunc(path)
	}
	return "Mock DOCX content extracted", nil
}

func TestMockDocConverter_Success(t *testing.T) {
	mockConverter := &MockDocConverter{
		BaseConverter: NewBaseConverter([]string{".docx", ".doc"},
			[]string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}),
		loadFunc: func(path string) (string, error) {
			return "# Document Title\n\nThis is extracted text from DOCX: " + path, nil
		},
	}

	result, err := mockConverter.Load("test.docx")
	if err != nil {
		t.Errorf("Mock Load() returned unexpected error: %v", err)
	}

	expected := "# Document Title\n\nThis is extracted text from DOCX: test.docx"
	if result != expected {
		t.Errorf("Mock Load() = %v, want %v", result, expected)
	}
}

func TestMockDocConverter_Error(t *testing.T) {
	mockConverter := &MockDocConverter{
		BaseConverter: NewBaseConverter([]string{".docx", ".doc"},
			[]string{"application/vnd.openxmlformats-officedocument.wordprocessingml.document"}),
		loadFunc: func(path string) (string, error) {
			return "", os.ErrNotExist
		},
	}

	_, err := mockConverter.Load("nonexistent.docx")
	if err == nil {
		t.Errorf("Mock Load() should return error")
	}
}

// Test various error scenarios
func TestDocConverter_LoadingScenarios(t *testing.T) {
	converter := NewDocConverter()

	testCases := []struct {
		name          string
		filename      string
		shouldError   bool
		errorContains string
	}{
		{
			name:          "Non-existent file",
			filename:      "/tmp/nonexistent-file-12345.docx",
			shouldError:   true,
			errorContains: "failed to convert document",
		},
		{
			name:          "Empty filename",
			filename:      "",
			shouldError:   true,
			errorContains: "failed to convert document",
		},
		{
			name:          "Directory instead of file",
			filename:      "/tmp",
			shouldError:   true,
			errorContains: "failed to convert document",
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

// Test the relationship and numbering structs for basic XML unmarshaling
func TestRelationshipStruct(t *testing.T) {
	// Test basic structure of Relationship
	rel := Relationship{
		ID:     "rId1",
		Type:   "http://example.com",
		Target: "target.xml",
	}

	if rel.ID != "rId1" {
		t.Errorf("Relationship.ID = %v, want rId1", rel.ID)
	}
	if rel.Type != "http://example.com" {
		t.Errorf("Relationship.Type = %v, want http://example.com", rel.Type)
	}
	if rel.Target != "target.xml" {
		t.Errorf("Relationship.Target = %v, want target.xml", rel.Target)
	}
}

func TestTextValStruct(t *testing.T) {
	// Test basic structure of TextVal
	tv := TextVal{
		Text: "sample text",
		Val:  "sample value",
	}

	if tv.Text != "sample text" {
		t.Errorf("TextVal.Text = %v, want 'sample text'", tv.Text)
	}
	if tv.Val != "sample value" {
		t.Errorf("TextVal.Val = %v, want 'sample value'", tv.Val)
	}
}
