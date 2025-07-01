package loaders

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPdfLoader_CanLoadMimeType(t *testing.T) {
	loader := &PdfLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "application/pdf",
			mimeType: "application/pdf",
			expected: true,
		},
		{
			name:     "application/pdf with charset",
			mimeType: "application/pdf; charset=utf-8",
			expected: true,
		},
		{
			name:     "text/plain",
			mimeType: "text/plain",
			expected: false,
		},
		{
			name:     "application/json",
			mimeType: "application/json",
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

const expectedPdfOpenError = "unable to open PDF file"

func TestPdfLoader_Load_FileNotFound(t *testing.T) {
	loader := &PdfLoader{}
	_, err := loader.Load("/nonexistent/file.pdf")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	if !strings.Contains(err.Error(), expectedPdfOpenError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedPdfOpenError)
	}
}

func TestPdfLoader_Load_InvalidPDF(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.pdf")

	// Create a file that's not a valid PDF
	invalidContent := "This is not a PDF file"
	err := os.WriteFile(testFile, []byte(invalidContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &PdfLoader{}
	_, err = loader.Load(testFile)
	if err == nil {
		t.Error("Load() should return error for invalid PDF file")
	}

	if !strings.Contains(err.Error(), expectedPdfOpenError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedPdfOpenError)
	}
}

func TestPdfLoader_Load_EmptyPDF(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.pdf")

	// Create an empty file
	err := os.WriteFile(testFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &PdfLoader{}
	_, err = loader.Load(testFile)
	if err == nil {
		t.Error("Load() should return error for empty file")
	}

	if !strings.Contains(err.Error(), expectedPdfOpenError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedPdfOpenError)
	}
}

func TestReadPdfFile_FileNotFound(t *testing.T) {
	_, err := readPdfFile("/nonexistent/file.pdf")
	if err == nil {
		t.Error("readPdfFile() should return error for nonexistent file")
	}

	if !strings.Contains(err.Error(), expectedPdfOpenError) {
		t.Errorf("readPdfFile() error = %v, should contain %v", err, expectedPdfOpenError)
	}
}

func TestReadPdfFile_InvalidPDF(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.pdf")

	// Create a file with minimal PDF header but invalid structure
	invalidPdf := "%PDF-1.4\n%%EOF"
	err := os.WriteFile(testFile, []byte(invalidPdf), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = readPdfFile(testFile)
	if err == nil {
		t.Error("readPdfFile() should return error for invalid PDF structure")
	}

	// Should contain error about opening the PDF
	if !strings.Contains(err.Error(), expectedPdfOpenError) {
		t.Errorf("readPdfFile() error = %v, should contain 'unable to open PDF file'", err)
	}
}

// Note: Testing with actual valid PDF files would require creating or embedding
// valid PDF content, which is complex. In a real-world scenario, you might:
// 1. Include small test PDF files in a testdata directory
// 2. Use a PDF generation library to create test PDFs
// 3. Mock the pdf.Open function for unit testing
//
// For now, we focus on testing error conditions and the CanLoadMimeType functionality.
