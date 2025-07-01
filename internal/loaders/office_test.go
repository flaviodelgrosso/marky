package loaders

import (
	"strings"
	"testing"
)

func TestDocLoader_CanLoadMimeType(t *testing.T) {
	loader := &DocLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			mimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			expected: true,
		},
		{
			name:     "application/vnd.openxmlformats-officedocument.wordprocessingml",
			mimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml",
			expected: true,
		},
		{
			name:     "DOCX with parameters",
			mimeType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document; charset=utf-8",
			expected: true,
		},
		{
			name:     "text/plain",
			mimeType: "text/plain",
			expected: false,
		},
		{
			name:     "application/pdf",
			mimeType: "application/pdf",
			expected: false,
		},
		{
			name:     "empty string",
			mimeType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loader.CanLoadMimeType(tt.mimeType)
			if result != tt.expected {
				t.Errorf("CanLoadMimeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDocLoader_Load_FileNotFound(t *testing.T) {
	loader := &DocLoader{}
	_, err := loader.Load("/nonexistent/file.docx")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	// Should contain error about reading or opening the file
	if !strings.Contains(err.Error(), "failed to") {
		t.Errorf("Load() error should contain 'failed to', got: %v", err)
	}
}

func TestExcelLoader_CanLoadMimeType(t *testing.T) {
	loader := &ExcelLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			mimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			expected: true,
		},
		{
			name:     "application/vnd.openxmlformats-officedocument.spreadsheetml",
			mimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml",
			expected: true,
		},
		{
			name:     "application/vnd.ms-excel",
			mimeType: "application/vnd.ms-excel",
			expected: true,
		},
		{
			name:     "Excel with parameters",
			mimeType: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet; charset=utf-8",
			expected: true,
		},
		{
			name:     "text/plain",
			mimeType: "text/plain",
			expected: false,
		},
		{
			name:     "application/pdf",
			mimeType: "application/pdf",
			expected: false,
		},
		{
			name:     "empty string",
			mimeType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loader.CanLoadMimeType(tt.mimeType)
			if result != tt.expected {
				t.Errorf("CanLoadMimeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExcelLoader_Load_FileNotFound(t *testing.T) {
	loader := &ExcelLoader{}
	_, err := loader.Load("/nonexistent/file.xlsx")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	// Should contain error about opening the file
	if !strings.Contains(err.Error(), "failed to") {
		t.Errorf("Load() error should contain 'failed to', got: %v", err)
	}
}

func TestPptxLoader_CanLoadMimeType(t *testing.T) {
	loader := &PptxLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			mimeType: "application/vnd.openxmlformats-officedocument.presentationml.presentation",
			expected: true,
		},
		{
			name:     "application/vnd.openxmlformats-officedocument.presentationml",
			mimeType: "application/vnd.openxmlformats-officedocument.presentationml",
			expected: true,
		},
		{
			name:     "PPTX with parameters",
			mimeType: "application/vnd.openxmlformats-officedocument.presentationml.presentation; charset=utf-8",
			expected: true,
		},
		{
			name:     "text/plain",
			mimeType: "text/plain",
			expected: false,
		},
		{
			name:     "application/pdf",
			mimeType: "application/pdf",
			expected: false,
		},
		{
			name:     "empty string",
			mimeType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loader.CanLoadMimeType(tt.mimeType)
			if result != tt.expected {
				t.Errorf("CanLoadMimeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPptxLoader_Load_FileNotFound(t *testing.T) {
	loader := &PptxLoader{}
	_, err := loader.Load("/nonexistent/file.pptx")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	// Should contain error about reading or opening the file
	if !strings.Contains(err.Error(), "failed to") {
		t.Errorf("Load() error should contain 'failed to', got: %v", err)
	}
}

// Test that all loaders implement the DocumentLoader interface
func TestDocumentLoaderInterface(t *testing.T) {
	loaders := []DocumentLoader{
		&CsvLoader{},
		&DocLoader{},
		&ExcelLoader{},
		&HTMLLoader{},
		&PdfLoader{},
		&PptxLoader{},
	}

	for i, loader := range loaders {
		t.Run(string(rune('A'+i)), func(t *testing.T) {
			// Test that CanLoadMimeType can be called
			_ = loader.CanLoadMimeType("test/type")

			// Test that Load can be called (will error for invalid paths, but should not panic)
			_, _ = loader.Load("/invalid/path")
		})
	}
}
