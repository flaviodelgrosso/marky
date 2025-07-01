package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile_Success(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	expectedContent := []byte("Hello, World!")

	err := os.WriteFile(testFile, expectedContent, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v, wantErr false", err)
		return
	}

	if string(result) != string(expectedContent) {
		t.Errorf("ReadFile() = %v, want %v", string(result), string(expectedContent))
	}
}

func TestReadFile_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.txt")

	err := os.WriteFile(testFile, []byte{}, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v, wantErr false", err)
		return
	}

	if len(result) != 0 {
		t.Errorf("ReadFile() = %v, want empty slice", result)
	}
}

func TestReadFile_FileNotFound(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("ReadFile() should return error for nonexistent file")
	}

	expectedError := "unable to read file"
	if err.Error()[:len(expectedError)] != expectedError {
		t.Errorf("ReadFile() error = %v, should start with %v", err, expectedError)
	}
}

func TestReadFile_LargeFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "large.txt")

	// Create a large file content
	largeContent := make([]byte, 10000)
	for i := range largeContent {
		largeContent[i] = byte('A' + (i % 26))
	}

	err := os.WriteFile(testFile, largeContent, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v, wantErr false", err)
		return
	}

	if len(result) != len(largeContent) {
		t.Errorf("ReadFile() length = %v, want %v", len(result), len(largeContent))
	}

	if string(result) != string(largeContent) {
		t.Error("ReadFile() content does not match expected large content")
	}
}

func TestReadFile_BinaryFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "binary.bin")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}

	err := os.WriteFile(testFile, binaryContent, 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v, wantErr false", err)
		return
	}

	if len(result) != len(binaryContent) {
		t.Errorf("ReadFile() length = %v, want %v", len(result), len(binaryContent))
	}

	for i, b := range result {
		if b != binaryContent[i] {
			t.Errorf("ReadFile() byte %d = %v, want %v", i, b, binaryContent[i])
		}
	}
}
