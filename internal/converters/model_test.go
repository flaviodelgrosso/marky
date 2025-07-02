package converters

import (
	"reflect"
	"testing"
)

func TestNewBaseConverter(t *testing.T) {
	extensions := []string{".txt", ".md"}
	mimeTypes := []string{"text/plain", "text/markdown"}

	converter := NewBaseConverter(extensions, mimeTypes)

	if !reflect.DeepEqual(converter.acceptedExtensions, extensions) {
		t.Errorf("NewBaseConverter() extensions = %v, want %v", converter.acceptedExtensions, extensions)
	}

	if !reflect.DeepEqual(converter.acceptedMimeTypes, mimeTypes) {
		t.Errorf("NewBaseConverter() mimeTypes = %v, want %v", converter.acceptedMimeTypes, mimeTypes)
	}
}

func TestBaseConverter_AcceptedExtensions(t *testing.T) {
	extensions := []string{".txt", ".md", ".html"}
	mimeTypes := []string{"text/plain", "text/markdown", "text/html"}

	converter := NewBaseConverter(extensions, mimeTypes)

	result := converter.AcceptedExtensions()
	if !reflect.DeepEqual(result, extensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", result, extensions)
	}
}

func TestBaseConverter_AcceptedMimeTypes(t *testing.T) {
	extensions := []string{".txt", ".md", ".html"}
	mimeTypes := []string{"text/plain", "text/markdown", "text/html"}

	converter := NewBaseConverter(extensions, mimeTypes)

	result := converter.AcceptedMimeTypes()
	if !reflect.DeepEqual(result, mimeTypes) {
		t.Errorf("AcceptedMimeTypes() = %v, want %v", result, mimeTypes)
	}
}

func TestBaseConverter_EmptyInputs(t *testing.T) {
	converter := NewBaseConverter([]string{}, []string{})

	if len(converter.AcceptedExtensions()) != 0 {
		t.Errorf("AcceptedExtensions() with empty input should return empty slice")
	}

	if len(converter.AcceptedMimeTypes()) != 0 {
		t.Errorf("AcceptedMimeTypes() with empty input should return empty slice")
	}
}

func TestBaseConverter_NilInputs(t *testing.T) {
	converter := NewBaseConverter(nil, nil)

	extensions := converter.AcceptedExtensions()
	mimeTypes := converter.AcceptedMimeTypes()

	// In Go, nil slices are valid and behave like empty slices for most operations
	// The functions return nil slices, which is acceptable behavior
	if len(extensions) != 0 {
		t.Errorf("AcceptedExtensions() with nil input should return empty slice, got %v", extensions)
	}

	if len(mimeTypes) != 0 {
		t.Errorf("AcceptedMimeTypes() with nil input should return empty slice, got %v", mimeTypes)
	}
}

// MockConverter is a simple implementation of the Converter interface for testing
type MockConverter struct {
	BaseConverter
	loadFunc func(path string) (string, error)
}

func (m *MockConverter) Load(path string) (string, error) {
	if m.loadFunc != nil {
		return m.loadFunc(path)
	}
	return "mock content", nil
}

func TestConverterInterface(t *testing.T) {
	extensions := []string{".mock"}
	mimeTypes := []string{"application/mock"}

	mock := &MockConverter{
		BaseConverter: NewBaseConverter(extensions, mimeTypes),
		loadFunc: func(path string) (string, error) {
			return "test content from " + path, nil
		},
	}

	// Test that MockConverter implements Converter interface
	var _ Converter = mock

	// Test AcceptedExtensions
	if !reflect.DeepEqual(mock.AcceptedExtensions(), extensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", mock.AcceptedExtensions(), extensions)
	}

	// Test AcceptedMimeTypes
	if !reflect.DeepEqual(mock.AcceptedMimeTypes(), mimeTypes) {
		t.Errorf("AcceptedMimeTypes() = %v, want %v", mock.AcceptedMimeTypes(), mimeTypes)
	}

	// Test Load
	result, err := mock.Load("test.mock")
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}
	expected := "test content from test.mock"
	if result != expected {
		t.Errorf("Load() = %v, want %v", result, expected)
	}
}
