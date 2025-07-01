package marky

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/flaviodelgrosso/marky/internal/loaders"
)

// MockLoader is a test double for testing purposes
type MockLoader struct {
	SupportedMimeTypes []string
	LoadResponse       string
	LoadError          error
}

func (m *MockLoader) Load(path string) (string, error) {
	if m.LoadError != nil {
		return "", m.LoadError
	}
	return m.LoadResponse, nil
}

func (m *MockLoader) CanLoadMimeType(mimeType string) bool {
	for _, supported := range m.SupportedMimeTypes {
		if strings.HasPrefix(mimeType, supported) {
			return true
		}
	}
	return false
}

func TestMarky_RegisterLoader(t *testing.T) {
	tests := []struct {
		name         string
		initialCount int
		addLoaders   int
	}{
		{
			name:         "Register single loader",
			initialCount: 0,
			addLoaders:   1,
		},
		{
			name:         "Register multiple loaders",
			initialCount: 2,
			addLoaders:   3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Marky{
				Loaders: make([]loaders.DocumentLoader, tt.initialCount),
			}

			for i := 0; i < tt.addLoaders; i++ {
				loader := &MockLoader{
					SupportedMimeTypes: []string{fmt.Sprintf("test/type-%d", i)},
				}
				m.RegisterLoader(loader)
			}

			expectedCount := tt.initialCount + tt.addLoaders
			if len(m.Loaders) != expectedCount {
				t.Errorf("Expected %d loaders, got %d", expectedCount, len(m.Loaders))
			}
		})
	}
}

func TestMarky_Convert_Success(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name           string
		filePath       string
		mimeType       string
		expectedResult string
		setupLoader    func() *MockLoader
	}{
		{
			name:           "Successful conversion with matching loader",
			filePath:       testFile,
			mimeType:       "text/plain",
			expectedResult: "converted content",
			setupLoader: func() *MockLoader {
				return &MockLoader{
					SupportedMimeTypes: []string{"text/plain"},
					LoadResponse:       "converted content",
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Marky{
				Loaders: []loaders.DocumentLoader{tt.setupLoader()},
			}

			result, err := m.Convert(tt.filePath)
			if err != nil {
				t.Errorf("Convert() error = %v, wantErr false", err)
				return
			}

			if result != tt.expectedResult {
				t.Errorf("Convert() = %v, want %v", result, tt.expectedResult)
			}
		})
	}
}

func TestMarky_Convert_NoMatchingLoader(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	m := &Marky{
		Loaders: []loaders.DocumentLoader{
			&MockLoader{
				SupportedMimeTypes: []string{"application/pdf"},
				LoadResponse:       "converted content",
			},
		},
	}

	_, err = m.Convert(testFile)
	if err == nil {
		t.Error("Convert() should return error when no matching loader is found")
	}

	expectedError := "no loader found for MIME type"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Convert() error = %v, should contain %v", err, expectedError)
	}
}

func TestMarky_Convert_LoaderError(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	expectedLoaderError := errors.New("loader failed")
	m := &Marky{
		Loaders: []loaders.DocumentLoader{
			&MockLoader{
				SupportedMimeTypes: []string{"text/plain"},
				LoadError:          expectedLoaderError,
			},
		},
	}

	_, err = m.Convert(testFile)
	if err == nil {
		t.Error("Convert() should return error when loader fails")
	}

	if !strings.Contains(err.Error(), expectedLoaderError.Error()) {
		t.Errorf("Convert() error = %v, should contain %v", err, expectedLoaderError)
	}
}

func TestMarky_Convert_FileNotFound(t *testing.T) {
	m := &Marky{
		Loaders: []loaders.DocumentLoader{},
	}

	_, err := m.Convert("/nonexistent/file.txt")
	if err == nil {
		t.Error("Convert() should return error for nonexistent file")
	}

	expectedError := "failed to detect MIME type"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Convert() error = %v, should contain %v", err, expectedError)
	}
}

func TestMarky_Convert_MultipleLoaders(t *testing.T) {
	// Create a temporary test file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")

	err := os.WriteFile(testFile, []byte("test content"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Create multiple loaders, first one should be used
	loader1 := &MockLoader{
		SupportedMimeTypes: []string{"text/plain"},
		LoadResponse:       "first loader result",
	}
	loader2 := &MockLoader{
		SupportedMimeTypes: []string{"text/plain"},
		LoadResponse:       "second loader result",
	}

	m := &Marky{
		Loaders: []loaders.DocumentLoader{loader1, loader2},
	}

	result, err := m.Convert(testFile)
	if err != nil {
		t.Errorf("Convert() error = %v, wantErr false", err)
		return
	}

	// Should use the first matching loader
	if result != "first loader result" {
		t.Errorf("Convert() = %v, want %v", result, "first loader result")
	}
}
