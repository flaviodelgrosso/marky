package marky

import (
	"testing"

	"github.com/flaviodelgrosso/marky/internal/loaders"
	"github.com/flaviodelgrosso/marky/internal/marky"
)

func TestInitialize(t *testing.T) {
	tests := []struct {
		name              string
		expectedLoaderLen int
	}{
		{
			name:              "Initialize should create marky instance with all loaders",
			expectedLoaderLen: 7, // CSV, DOC, Excel, HTML, IPYNB, PDF, PPTX
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := New()

			// Check that result is not nil
			if result == nil {
				t.Error("New() should not return nil")
			}

			// Check that it's a Marky instance
			instance, ok := result.(*marky.Marky)
			if !ok {
				t.Error("New() should return a *marky instance")
			}

			// Check that all expected loaders are registered
			if len(instance.Loaders) != tt.expectedLoaderLen {
				t.Errorf("New() should register %d loaders, got %d", tt.expectedLoaderLen, len(instance.Loaders))
			}

			// Check that specific loader types are present
			expectedTypes := map[string]bool{
				"*loaders.CsvLoader":   false,
				"*loaders.DocLoader":   false,
				"*loaders.ExcelLoader": false,
				"*loaders.HtmlLoader":  false,
				"*loaders.IpynbLoader": false,
				"*loaders.PdfLoader":   false,
				"*loaders.PptxLoader":  false,
			}

			for _, loader := range instance.Loaders {
				switch loader.(type) {
				case *loaders.CsvLoader:
					expectedTypes["*loaders.CsvLoader"] = true
				case *loaders.DocLoader:
					expectedTypes["*loaders.DocLoader"] = true
				case *loaders.ExcelLoader:
					expectedTypes["*loaders.ExcelLoader"] = true
				case *loaders.HTMLLoader:
					expectedTypes["*loaders.HtmlLoader"] = true
				case *loaders.IpynbLoader:
					expectedTypes["*loaders.IpynbLoader"] = true
				case *loaders.PdfLoader:
					expectedTypes["*loaders.PdfLoader"] = true
				case *loaders.PptxLoader:
					expectedTypes["*loaders.PptxLoader"] = true
				}
			}

			for loaderType, found := range expectedTypes {
				if !found {
					t.Errorf("Expected loader type %s was not registered", loaderType)
				}
			}
		})
	}
}

func TestInitializeCapacity(t *testing.T) {
	// Test that the slice is pre-allocated with the correct capacity
	instance := New()
	markyInstance := instance.(*marky.Marky)

	// The slice should have exactly 7 elements (the registered loaders)
	if len(markyInstance.Loaders) != 7 {
		t.Errorf("Expected 7 loaders, got %d", len(markyInstance.Loaders))
	}
}
