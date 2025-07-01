package mimetypes

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDetectMimeType_Success(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		fileName     string
		fileContent  []byte
		expectedMime string
		expectedExt  string
	}{
		{
			name:         "Text file",
			fileName:     "test.txt",
			fileContent:  []byte("Hello, World!"),
			expectedMime: "text/plain", // Accept both variants
			expectedExt:  ".txt",
		},
		{
			name:         "HTML file",
			fileName:     "test.html",
			fileContent:  []byte("<!DOCTYPE html><html><body>Hello</body></html>"),
			expectedMime: "text/html",
			expectedExt:  ".html",
		},
		{
			name:         "CSV file",
			fileName:     "test.csv",
			fileContent:  []byte("name,age\nJohn,30\nJane,25"),
			expectedMime: "text/", // CSV can be detected as text/csv or text/plain
			expectedExt:  ".csv",
		},
		{
			name:         "PDF file",
			fileName:     "test.pdf",
			fileContent:  []byte("%PDF-1.4\n%âãÏÓ"),
			expectedMime: "application/pdf",
			expectedExt:  ".pdf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, tt.fileContent, 0o644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := DetectMimeType(testFile)
			if err != nil {
				t.Errorf("DetectMimeType() error = %v, wantErr false", err)
				return
			}

			if !strings.HasPrefix(result.MimeType, tt.expectedMime) {
				t.Errorf("DetectMimeType() MimeType = %v, want %v (or variant)", result.MimeType, tt.expectedMime)
			}

			if result.Extension != tt.expectedExt {
				t.Errorf("DetectMimeType() Extension = %v, want %v", result.Extension, tt.expectedExt)
			}
		})
	}
}

func TestDetectMimeType_FileNotFound(t *testing.T) {
	_, err := DetectMimeType("/nonexistent/file.txt")
	if err == nil {
		t.Error("DetectMimeType() should return error for nonexistent file")
	}
}

func TestDetectMimeType_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.txt")

	err := os.WriteFile(testFile, []byte{}, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := DetectMimeType(testFile)
	if err != nil {
		t.Errorf("DetectMimeType() error = %v, wantErr false", err)
		return
	}

	if result.Extension != ".txt" {
		t.Errorf("DetectMimeType() Extension = %v, want %v", result.Extension, ".txt")
	}

	// Empty files are typically detected as text/plain or application/octet-stream
	if !strings.HasPrefix(result.MimeType, "text/plain") &&
		!strings.HasPrefix(result.MimeType, "application/octet-stream") {
		t.Errorf(
			"DetectMimeType() MimeType = %v, want text/plain or application/octet-stream (or variant)",
			result.MimeType,
		)
	}
}

func TestDetectMimeType_LargeFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large.txt")

	// Create a file larger than 512 bytes to test the sampling logic
	largeContent := make([]byte, 1024)
	for i := range largeContent {
		largeContent[i] = byte('A' + (i % 26))
	}

	err := os.WriteFile(testFile, largeContent, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := DetectMimeType(testFile)
	if err != nil {
		t.Errorf("DetectMimeType() error = %v, wantErr false", err)
		return
	}

	if result.Extension != ".txt" {
		t.Errorf("DetectMimeType() Extension = %v, want %v", result.Extension, ".txt")
	}

	if !strings.HasPrefix(result.MimeType, "text/plain") {
		t.Errorf("DetectMimeType() MimeType = %v, want text/plain (or variant)", result.MimeType)
	}
}

func TestIsMimeTypeSupported(t *testing.T) {
	tests := []struct {
		name           string
		mimeType       string
		supportedTypes []string
		expected       bool
	}{
		{
			name:           "Exact match",
			mimeType:       "text/csv",
			supportedTypes: []string{"text/csv", "application/pdf"},
			expected:       true,
		},
		{
			name:           "Prefix match",
			mimeType:       "text/csv; charset=utf-8",
			supportedTypes: []string{"text/csv", "application/pdf"},
			expected:       true,
		},
		{
			name:           "No match",
			mimeType:       "application/json",
			supportedTypes: []string{"text/csv", "application/pdf"},
			expected:       false,
		},
		{
			name:           "Empty supported types",
			mimeType:       "text/csv",
			supportedTypes: []string{},
			expected:       false,
		},
		{
			name:           "Empty mime type",
			mimeType:       "",
			supportedTypes: []string{"text/csv"},
			expected:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMimeTypeSupported(tt.mimeType, tt.supportedTypes)
			if result != tt.expected {
				t.Errorf("IsMimeTypeSupported() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetSupportedMimeTypes(t *testing.T) {
	result := GetSupportedMimeTypes()

	// Check that all expected file types are present
	expectedTypes := []string{"csv", "docx", "doc", "excel", "html", "pdf", "pptx"}
	for _, fileType := range expectedTypes {
		if _, exists := result[fileType]; !exists {
			t.Errorf("GetSupportedMimeTypes() missing file type: %s", fileType)
		}
	}

	// Check specific MIME types for known file types
	csvTypes := result["csv"]
	if len(csvTypes) == 0 {
		t.Error("GetSupportedMimeTypes() csv should have MIME types")
	}

	found := false
	for _, mimeType := range csvTypes {
		if mimeType == "text/csv" {
			found = true
			break
		}
	}
	if !found {
		t.Error("GetSupportedMimeTypes() csv should include 'text/csv'")
	}

	// Check PDF MIME types
	pdfTypes := result["pdf"]
	if len(pdfTypes) != 1 || pdfTypes[0] != "application/pdf" {
		t.Errorf("GetSupportedMimeTypes() pdf = %v, want [application/pdf]", pdfTypes)
	}
}

func TestEnhanceMimeTypeDetection_OfficeFormats(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name         string
		fileName     string
		fileContent  []byte
		baseMimeType string
		expectedMime string
	}{
		{
			name:         "ZIP signature with DOCX extension",
			fileName:     "test.docx",
			fileContent:  []byte{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00}, // ZIP signature
			baseMimeType: "application/zip",
			expectedMime: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		},
		{
			name:         "ZIP signature with XLSX extension",
			fileName:     "test.xlsx",
			fileContent:  []byte{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00}, // ZIP signature
			baseMimeType: "application/zip",
			expectedMime: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		},
		{
			name:         "ZIP signature with PPTX extension",
			fileName:     "test.pptx",
			fileContent:  []byte{0x50, 0x4B, 0x03, 0x04, 0x14, 0x00}, // ZIP signature
			baseMimeType: "application/zip",
			expectedMime: "application/vnd.openxmlformats-officedocument.presentationml.presentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, tt.fileName)
			err := os.WriteFile(testFile, tt.fileContent, 0o644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			result, err := DetectMimeType(testFile)
			if err != nil {
				t.Errorf("DetectMimeType() error = %v, wantErr false", err)
				return
			}

			if result.MimeType != tt.expectedMime {
				t.Errorf("DetectMimeType() MimeType = %v, want %v", result.MimeType, tt.expectedMime)
			}
		})
	}
}
