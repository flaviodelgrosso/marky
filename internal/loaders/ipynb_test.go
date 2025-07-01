package loaders

import (
	"os"
	"strings"
	"testing"
)

func TestIpynbLoader_Load(t *testing.T) {
	// Create a temporary test notebook file
	notebookContent := `{
		"nbformat": 4,
		"nbformat_minor": 4,
		"cells": [
			{
				"cell_type": "markdown",
				"source": [
					"# Test Notebook\n",
					"\n",
					"This is a test notebook for the IPYNB loader."
				]
			},
			{
				"cell_type": "code",
				"source": [
					"print('Hello, World!')\n",
					"x = 1 + 1\n",
					"print(f'1 + 1 = {x}')"
				]
			},
			{
				"cell_type": "markdown",
				"source": [
					"## Results\n",
					"\n",
					"The code above should print:\n",
					"- Hello, World!\n",
					"- 1 + 1 = 2"
				]
			},
			{
				"cell_type": "raw",
				"source": [
					"This is raw text\n",
					"that should be preserved as-is"
				]
			}
		],
		"metadata": {
			"title": "Test Notebook"
		}
	}`

	// Create temporary file
	tmpFile, err := os.CreateTemp(t.TempDir(), "test_notebook_*.ipynb")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(notebookContent); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tmpFile.Close()

	// Test the loader
	loader := &IpynbLoader{}
	result, err := loader.Load(tmpFile.Name())
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify the result contains expected content
	expectedContent := []string{
		"# Test Notebook",
		"This is a test notebook",
		"```python",
		"print('Hello, World!')",
		"x = 1 + 1",
		"print(f'1 + 1 = {x}')",
		"```",
		"## Results",
		"The code above should print:",
		"- Hello, World!",
		"- 1 + 1 = 2",
		"```",
		"This is raw text",
		"that should be preserved as-is",
		"```",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected content not found in result: %q", expected)
		}
	}
}

func TestIpynbLoader_CanLoadMimeType(t *testing.T) {
	loader := &IpynbLoader{}

	tests := []struct {
		mimeType string
		expected bool
	}{
		{"application/json", true},
		{"application/x-ipynb+json", true},
		{"text/plain", false},
		{"application/pdf", false},
		{"", false},
	}

	for _, test := range tests {
		result := loader.CanLoadMimeType(test.mimeType)
		if result != test.expected {
			t.Errorf("CanLoadMimeType(%q) = %v, expected %v", test.mimeType, result, test.expected)
		}
	}
}

func TestConvertNotebookToMarkdown(t *testing.T) {
	notebook := JupyterNotebook{
		NBFormat:      4,
		NBFormatMinor: 4,
		Cells: []NotebookCell{
			{
				CellType: "markdown",
				Source:   []string{"# Main Title\n", "This is markdown content."},
			},
			{
				CellType: "code",
				Source:   []string{"print('test')\n", "x = 42"},
			},
			{
				CellType: "raw",
				Source:   []string{"Raw content\n", "line 2"},
			},
		},
		Metadata: NotebookMetadata{
			Title: "Notebook Title",
		},
	}

	result := convertNotebookToMarkdown(notebook)

	expectedContent := []string{
		"# Main Title",
		"This is markdown content.",
		"```python",
		"print('test')",
		"x = 42",
		"```",
		"```",
		"Raw content",
		"line 2",
		"```",
	}

	for _, expected := range expectedContent {
		if !strings.Contains(result, expected) {
			t.Errorf("Expected content not found in result: %q", expected)
		}
	}
}

func TestConvertNotebookToMarkdown_EmptyTitle(t *testing.T) {
	notebook := JupyterNotebook{
		NBFormat:      4,
		NBFormatMinor: 4,
		Cells: []NotebookCell{
			{
				CellType: "markdown",
				Source:   []string{"Some content without title."},
			},
		},
		Metadata: NotebookMetadata{
			Title: "Metadata Title",
		},
	}

	result := convertNotebookToMarkdown(notebook)

	// Should add title from metadata when no # heading is found in cells
	if !strings.Contains(result, "# Metadata Title") {
		t.Error("Expected title from metadata to be added")
	}
}
