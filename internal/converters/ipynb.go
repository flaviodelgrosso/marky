package converters

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// IpynbConverter handles loading and converting Jupyter Notebook (.ipynb) files to markdown.
type IpynbConverter struct {
	BaseConverter
}

// NewIpynbConverter creates a new Jupyter Notebook converter with appropriate MIME types and extensions.
func NewIpynbConverter() Converter {
	return &IpynbConverter{
		BaseConverter: NewBaseConverter(
			[]string{".ipynb"},
			[]string{"application/x-ipynb+json", "application/json"},
		),
	}
}

// NotebookCell represents a cell in a Jupyter notebook.
type NotebookCell struct {
	CellType string   `json:"cell_type"`
	Source   []string `json:"source"`
}

// NotebookMetadata represents the metadata section of a Jupyter notebook.
type NotebookMetadata struct {
	Title string `json:"title,omitempty"`
}

// JupyterNotebook represents the structure of a Jupyter notebook file.
type JupyterNotebook struct {
	NBFormat      int              `json:"nbformat"`
	NBFormatMinor int              `json:"nbformat_minor"`
	Cells         []NotebookCell   `json:"cells"`
	Metadata      NotebookMetadata `json:"metadata"`
}

// Load reads a Jupyter notebook file and converts it to markdown.
func (*IpynbConverter) Load(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read ipynb file: %w", err)
	}

	// Parse the JSON content
	var notebook JupyterNotebook
	if err := json.Unmarshal(content, &notebook); err != nil {
		return "", fmt.Errorf("failed to parse ipynb file: %w", err)
	}

	return convertNotebookToMarkdown(notebook), nil
}

// convertNotebookToMarkdown converts a Jupyter notebook to markdown format.
func convertNotebookToMarkdown(notebook JupyterNotebook) string {
	var mdParts []string
	var title string

	for _, cell := range notebook.Cells {
		cellContent := strings.Join(cell.Source, "")

		switch cell.CellType {
		case "markdown":
			mdParts = append(mdParts, cellContent)

			// Extract the first # heading as title if not already found
			if title == "" {
				lines := strings.SplitSeq(cellContent, "\n")
				for line := range lines {
					trimmed := strings.TrimSpace(line)
					if after, ok := strings.CutPrefix(trimmed, "# "); ok {
						title = strings.TrimSpace(after)
						break
					}
				}
			}

		case "code":
			// Code cells are wrapped in Markdown code blocks
			if strings.TrimSpace(cellContent) != "" {
				mdParts = append(mdParts, fmt.Sprintf("```python\n%s\n```", cellContent))
			}

		case "raw":
			// Raw cells are wrapped in plain code blocks
			if strings.TrimSpace(cellContent) != "" {
				mdParts = append(mdParts, fmt.Sprintf("```\n%s\n```", cellContent))
			}
		}
	}

	// Check for title in notebook metadata if not found in cells
	if title == "" && notebook.Metadata.Title != "" {
		title = notebook.Metadata.Title
	}

	// Add title as first heading if found
	markdown := strings.Join(mdParts, "\n\n")
	if title != "" && !strings.HasPrefix(strings.TrimSpace(markdown), "# ") {
		markdown = fmt.Sprintf("# %s\n\n%s", title, markdown)
	}

	return markdown
}
