package utils

import (
	"bytes"
	"fmt"
	"strings"
)

// ToMarkdownTable converts a 2D string slice to a markdown table format.
func ToMarkdownTable(rows [][]string) string {
	if len(rows) == 0 {
		return ""
	}

	// Additional safety check for empty first row
	if len(rows[0]) == 0 {
		return ""
	}

	var buf bytes.Buffer
	headerColCount := len(rows[0])

	// Header
	buf.WriteString("|")
	for _, cell := range rows[0] {
		// Escape pipe characters in cell content and trim whitespace
		escapedCell := strings.ReplaceAll(strings.TrimSpace(cell), "|", "\\|")
		fmt.Fprintf(&buf, " %s |", escapedCell)
	}
	buf.WriteString("\n|")

	// Header separator
	for range headerColCount {
		buf.WriteString(" --- |")
	}
	buf.WriteString("\n")

	// Data rows - only process if we have more than one row
	if len(rows) > 1 {
		for _, row := range rows[1:] {
			buf.WriteString("|")
			// Handle rows with different column counts
			for i := range headerColCount {
				var cell string
				if i < len(row) {
					// Escape pipe characters in cell content and trim whitespace
					cell = strings.ReplaceAll(strings.TrimSpace(row[i]), "|", "\\|")
				}
				fmt.Fprintf(&buf, " %s |", cell)
			}
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
