package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewIpynbConverter(t *testing.T) {
	converter := NewIpynbConverter()

	expectedExtensions := []string{".ipynb"}
	expectedMimeTypes := []string{"application/x-ipynb+json", "application/json"}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewIpynbConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewIpynbConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestIpynbConverter_Load_ValidNotebook(t *testing.T) {
	// Create a temporary Jupyter notebook file
	tempDir := t.TempDir()
	ipynbFile := filepath.Join(tempDir, "test.ipynb")

	notebookContent := `{
 "cells": [
  {
   "cell_type": "markdown",
   "source": [
    "# Test Notebook\n",
    "\n",
    "This is a test notebook."
   ]
  },
  {
   "cell_type": "code",
   "source": [
    "print('Hello, World!')\n",
    "x = 42"
   ]
  }
 ],
 "metadata": {
  "title": "Test Notebook"
 },
 "nbformat": 4,
 "nbformat_minor": 4
}`

	err := os.WriteFile(ipynbFile, []byte(notebookContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test notebook file: %v", err)
	}

	converter := NewIpynbConverter()
	result, err := converter.Load(ipynbFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that markdown cells are preserved
	if !strings.Contains(result, "# Test Notebook") {
		t.Errorf("Load() should preserve markdown content")
	}

	// Check that code cells are converted to code blocks
	if !strings.Contains(result, "```") {
		t.Errorf("Load() should convert code cells to code blocks")
	}

	if !strings.Contains(result, "print('Hello, World!')") {
		t.Errorf("Load() should preserve code content")
	}
}

func TestIpynbConverter_Load_EmptyNotebook(t *testing.T) {
	// Create a temporary empty notebook file
	tempDir := t.TempDir()
	ipynbFile := filepath.Join(tempDir, "empty.ipynb")

	notebookContent := `{
 "cells": [],
 "metadata": {},
 "nbformat": 4,
 "nbformat_minor": 4
}`

	err := os.WriteFile(ipynbFile, []byte(notebookContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test notebook file: %v", err)
	}

	converter := NewIpynbConverter()
	result, err := converter.Load(ipynbFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Empty notebook should return minimal content
	if len(result) > 10 { // Allow for some minimal formatting
		t.Errorf("Load() with empty notebook should return minimal content, got: %v", result)
	}
}

func TestIpynbConverter_Load_OnlyMarkdown(t *testing.T) {
	// Create a temporary notebook file with only markdown cells
	tempDir := t.TempDir()
	ipynbFile := filepath.Join(tempDir, "markdown_only.ipynb")

	notebookContent := `{
 "cells": [
  {
   "cell_type": "markdown",
   "source": [
    "# Title\n",
    "\n",
    "Some text here."
   ]
  },
  {
   "cell_type": "markdown",
   "source": [
    "## Subtitle\n",
    "\n",
    "More text."
   ]
  }
 ],
 "metadata": {},
 "nbformat": 4,
 "nbformat_minor": 4
}`

	err := os.WriteFile(ipynbFile, []byte(notebookContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test notebook file: %v", err)
	}

	converter := NewIpynbConverter()
	result, err := converter.Load(ipynbFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that markdown content is preserved
	if !strings.Contains(result, "# Title") {
		t.Errorf("Load() should preserve markdown titles")
	}
	if !strings.Contains(result, "## Subtitle") {
		t.Errorf("Load() should preserve markdown subtitles")
	}
}

func TestIpynbConverter_Load_OnlyCode(t *testing.T) {
	// Create a temporary notebook file with only code cells
	tempDir := t.TempDir()
	ipynbFile := filepath.Join(tempDir, "code_only.ipynb")

	notebookContent := `{
 "cells": [
  {
   "cell_type": "code",
   "source": [
    "import numpy as np\n",
    "arr = np.array([1, 2, 3])"
   ]
  },
  {
   "cell_type": "code",
   "source": [
    "print(arr)\n",
    "print('Done')"
   ]
  }
 ],
 "metadata": {},
 "nbformat": 4,
 "nbformat_minor": 4
}`

	err := os.WriteFile(ipynbFile, []byte(notebookContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test notebook file: %v", err)
	}

	converter := NewIpynbConverter()
	result, err := converter.Load(ipynbFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that code content is converted to code blocks
	if !strings.Contains(result, "```") {
		t.Errorf("Load() should convert code cells to code blocks")
	}
	if !strings.Contains(result, "import numpy as np") {
		t.Errorf("Load() should preserve code content")
	}
}

func TestIpynbConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewIpynbConverter()
	_, err := converter.Load("/nonexistent/file.ipynb")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read ipynb file") {
		t.Errorf("Load() error should mention ipynb file reading failure")
	}
}

func TestIpynbConverter_Load_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.ipynb")

	invalidContent := `{
 "cells": [
  {
   "cell_type": "markdown",
   "source": "unclosed quote
  }
 ],
 "metadata": {},
 "nbformat": 4`

	err := os.WriteFile(invalidFile, []byte(invalidContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewIpynbConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid JSON")
	}

	if !strings.Contains(err.Error(), "failed to parse ipynb file") {
		t.Errorf("Load() error should mention ipynb parsing failure")
	}
}

func TestIpynbConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary file with no read permissions
	tempDir := t.TempDir()
	noPermFile := filepath.Join(tempDir, "noperm.ipynb")

	err := os.WriteFile(noPermFile, []byte(`{"cells":[],"metadata":{},"nbformat":4}`), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	converter := NewIpynbConverter()
	_, err = converter.Load(noPermFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestIpynbConverter_Interface(t *testing.T) {
	converter := NewIpynbConverter()

	// Test that IpynbConverter implements Converter interface

	// Test AcceptedExtensions
	extensions := converter.AcceptedExtensions()
	expectedExtensions := []string{".ipynb"}
	if !reflect.DeepEqual(extensions, expectedExtensions) {
		t.Errorf("AcceptedExtensions() = %v, want %v", extensions, expectedExtensions)
	}

	// Test AcceptedMimeTypes
	mimeTypes := converter.AcceptedMimeTypes()
	if len(mimeTypes) != 2 {
		t.Errorf("AcceptedMimeTypes() should return 2 MIME types, got %d", len(mimeTypes))
	}
}

// Test the Jupyter Notebook structs
func TestNotebookCellStruct(t *testing.T) {
	cell := NotebookCell{
		CellType: "markdown",
		Source:   []string{"# Title\n", "Some content"},
	}

	if cell.CellType != "markdown" {
		t.Errorf("NotebookCell.CellType = %v, want 'markdown'", cell.CellType)
	}

	if len(cell.Source) != 2 {
		t.Errorf("NotebookCell.Source should have 2 elements, got %d", len(cell.Source))
	}

	if cell.Source[0] != "# Title\n" {
		t.Errorf("NotebookCell.Source[0] = %v, want '# Title\\n'", cell.Source[0])
	}
}

func TestNotebookMetadataStruct(t *testing.T) {
	metadata := NotebookMetadata{
		Title: "Test Notebook",
	}

	if metadata.Title != "Test Notebook" {
		t.Errorf("NotebookMetadata.Title = %v, want 'Test Notebook'", metadata.Title)
	}
}

func TestJupyterNotebookStruct(t *testing.T) {
	notebook := JupyterNotebook{
		NBFormat:      4,
		NBFormatMinor: 4,
		Cells: []NotebookCell{
			{
				CellType: "code",
				Source:   []string{"print('hello')"},
			},
		},
		Metadata: NotebookMetadata{
			Title: "Test",
		},
	}

	if notebook.NBFormat != 4 {
		t.Errorf("JupyterNotebook.NBFormat = %v, want 4", notebook.NBFormat)
	}

	if notebook.NBFormatMinor != 4 {
		t.Errorf("JupyterNotebook.NBFormatMinor = %v, want 4", notebook.NBFormatMinor)
	}

	if len(notebook.Cells) != 1 {
		t.Errorf("JupyterNotebook.Cells should have 1 cell, got %d", len(notebook.Cells))
	}

	if notebook.Metadata.Title != "Test" {
		t.Errorf("JupyterNotebook.Metadata.Title = %v, want 'Test'", notebook.Metadata.Title)
	}
}

// Test various error scenarios
func TestIpynbConverter_LoadingScenarios(t *testing.T) {
	converter := NewIpynbConverter()

	testCases := []struct {
		name          string
		filename      string
		shouldError   bool
		errorContains string
	}{
		{
			name:          "Non-existent file",
			filename:      "/tmp/nonexistent-file-12345.ipynb",
			shouldError:   true,
			errorContains: "failed to read ipynb file",
		},
		{
			name:          "Empty filename",
			filename:      "",
			shouldError:   true,
			errorContains: "failed to read ipynb file",
		},
		{
			name:          "Directory instead of file",
			filename:      "/tmp",
			shouldError:   true,
			errorContains: "failed to read ipynb file",
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
