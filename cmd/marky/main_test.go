package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flaviodelgrosso/marky"
	markyInternal "github.com/flaviodelgrosso/marky/internal/marky"
)

func TestMain_Convert_CSV(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test CSV file
	csvContent := "Name,Age,City\nJohn,30,New York\nJane,25,Los Angeles"
	csvFile := filepath.Join(tempDir, "test.csv")
	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	// Test conversion
	md := marky.Initialize()
	result, err := md.Convert(csvFile)
	if err != nil {
		t.Errorf("Convert() error = %v, wantErr false", err)
		return
	}

	expectedContent := []string{
		"| Name | Age | City |",
		"| --- | --- | --- |",
		"| John | 30 | New York |",
		"| Jane | 25 | Los Angeles |",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Convert() result should contain %v", expected)
		}
	}
}

func TestMain_Convert_HTML(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test HTML file
	htmlContent := "<html><body><h1>Hello World</h1><p>This is a <strong>test</strong>.</p></body></html>"
	htmlFile := filepath.Join(tempDir, "test.html")
	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	// Test conversion
	md := marky.Initialize()
	result, err := md.Convert(htmlFile)
	if err != nil {
		t.Errorf("Convert() error = %v, wantErr false", err)
		return
	}

	expectedContent := []string{
		"# Hello World",
		"**test**",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Convert() result should contain %v", expected)
		}
	}
}

func TestMain_Convert_UnsupportedFormat(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file with unsupported format
	unsupportedFile := filepath.Join(tempDir, "test.unknown")
	err := os.WriteFile(unsupportedFile, []byte("some content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test conversion should fail
	md := marky.Initialize()
	_, err = md.Convert(unsupportedFile)
	if err == nil {
		t.Error("Convert() should return error for unsupported file format")
	}

	expectedError := "no loader found for MIME type"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Convert() error = %v, should contain %v", err, expectedError)
	}
}

func TestMain_Convert_NonexistentFile(t *testing.T) {
	// Test conversion of nonexistent file
	md := marky.Initialize()
	_, err := md.Convert("/nonexistent/file.txt")
	if err == nil {
		t.Error("Convert() should return error for nonexistent file")
	}

	expectedError := "failed to detect MIME type"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Convert() error = %v, should contain %v", err, expectedError)
	}
}

func TestMain_AllLoadersRegistered(t *testing.T) {
	md := marky.Initialize()

	// Cast to concrete type to access Loaders field
	markyInstance, ok := md.(*markyInternal.Marky)
	if !ok {
		t.Fatal("Initialize() should return *marky.marky")
	}

	// Check that we have the expected number of loaders
	expectedLoaderCount := 7 // CSV, DOC, Excel, HTML, IPYNB, PDF, PPTX
	if len(markyInstance.Loaders) != expectedLoaderCount {
		t.Errorf("Expected %d loaders, got %d", expectedLoaderCount, len(markyInstance.Loaders))
	}

	// Test that each loader can handle at least one MIME type
	for i, loader := range markyInstance.Loaders {
		found := false
		testMimeTypes := []string{
			"text/csv",
			"application/csv",
			"text/html",
			"application/pdf",
			"application/json",
			"application/x-ipynb+json",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			"application/vnd.openxmlformats-officedocument.presentationml.presentation",
		}

		for _, mimeType := range testMimeTypes {
			if loader.CanLoadMimeType(mimeType) {
				found = true
				break
			}
		}

		if !found {
			t.Errorf("Loader %d should handle at least one of the test MIME types", i)
		}
	}
}

// Integration test for file writing functionality (simulating CLI behavior)
func TestMain_WriteToFile(t *testing.T) {
	tempDir := t.TempDir()

	// Create a test CSV file
	csvContent := "Name,Value\nTest,123"
	csvFile := filepath.Join(tempDir, "test.csv")
	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	// Convert to markdown
	md := marky.Initialize()
	result, err := md.Convert(csvFile)
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	// Write to output file (simulating CLI --output functionality)
	outputFile := filepath.Join(tempDir, "output.md")
	err = os.WriteFile(outputFile, []byte(result), 0o644)
	if err != nil {
		t.Fatalf("Failed to write output file: %v", err)
	}

	// Read back and verify
	writtenContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if string(writtenContent) != result {
		t.Error("Written content should match conversion result")
	}

	// Verify it contains expected markdown table structure
	content := string(writtenContent)
	if !strings.Contains(content, "| Name | Value |") {
		t.Error("Output file should contain markdown table header")
	}
	if !strings.Contains(content, "| Test | 123 |") {
		t.Error("Output file should contain markdown table data")
	}
}
