package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewHTMLConverter(t *testing.T) {
	converter := NewHTMLConverter()

	expectedExtensions := []string{".html", ".htm"}
	expectedMimeTypes := []string{"text/html"}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewHTMLConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewHTMLConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestHTMLConverter_Load_ValidHTML(t *testing.T) {
	// Create a temporary HTML file
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "test.html")

	htmlContent := `<!DOCTYPE html>
<html>
<head>
    <title>Test Document</title>
</head>
<body>
    <h1>Hello World</h1>
    <p>This is a <strong>test</strong> paragraph.</p>
    <ul>
        <li>Item 1</li>
        <li>Item 2</li>
    </ul>
</body>
</html>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that HTML is converted to markdown
	if !strings.Contains(result, "Hello World") {
		t.Errorf("Load() should preserve text content")
	}
	if !strings.Contains(result, "test") {
		t.Errorf("Load() should preserve text content with formatting")
	}
	// Check for list items (may be converted differently by html-to-markdown)
	if !strings.Contains(result, "Item 1") || !strings.Contains(result, "Item 2") {
		t.Errorf("Load() should preserve list content")
	}
}

func TestHTMLConverter_Load_SimpleHTML(t *testing.T) {
	// Create a temporary HTML file with simple content
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "simple.html")

	htmlContent := `<h2>Title</h2>
<p>Simple paragraph with <em>emphasis</em>.</p>
<a href="https://example.com">Link</a>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check markdown conversion
	if !strings.Contains(result, "## Title") {
		t.Errorf("Load() should convert h2 to markdown heading")
	}
	if !strings.Contains(result, "*emphasis*") {
		t.Errorf("Load() should convert em to markdown italic")
	}
	if !strings.Contains(result, "[Link](https://example.com)") {
		t.Errorf("Load() should convert a to markdown link")
	}
}

func TestHTMLConverter_Load_EmptyHTML(t *testing.T) {
	// Create a temporary empty HTML file
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "empty.html")

	err := os.WriteFile(htmlFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Empty HTML should return empty or minimal markdown
	if len(result) > 10 { // Allow for some whitespace/newlines
		t.Errorf("Load() with empty HTML should return minimal content, got: %v", result)
	}
}

func TestHTMLConverter_Load_HTMLWithTable(t *testing.T) {
	// Create a temporary HTML file with a table
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "table.html")

	htmlContent := `<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Age</th>
        </tr>
    </thead>
    <tbody>
        <tr>
            <td>John</td>
            <td>30</td>
        </tr>
        <tr>
            <td>Jane</td>
            <td>25</td>
        </tr>
    </tbody>
</table>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that table is converted to markdown format (may use different formatting)
	if !strings.Contains(result, "Name") && !strings.Contains(result, "John") {
		t.Errorf("Load() should preserve table content")
	}
	if !strings.Contains(result, "Age") && !strings.Contains(result, "30") {
		t.Errorf("Load() should preserve table data")
	}
}

func TestHTMLConverter_Load_HTMLWithSpecialCharacters(t *testing.T) {
	// Create a temporary HTML file with special characters
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "special.html")

	htmlContent := `<p>&lt;script&gt;alert('xss')&lt;/script&gt;</p>
<p>&amp; ampersand</p>
<p>&quot;quoted text&quot;</p>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that HTML entities are properly decoded (may be decoded differently)
	if !strings.Contains(result, "script") || !strings.Contains(result, "alert") {
		t.Errorf("Load() should decode and preserve script content")
	}
	if !strings.Contains(result, "ampersand") {
		t.Errorf("Load() should preserve ampersand text")
	}
}

func TestHTMLConverter_Load_HTMLWithCode(t *testing.T) {
	// Create a temporary HTML file with code elements
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "code.html")

	htmlContent := `<p>Inline <code>code block</code> here.</p>
<pre><code>
function hello() {
    console.log("Hello, world!");
}
</code></pre>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that code is converted to markdown code format
	if !strings.Contains(result, "`code block`") {
		t.Errorf("Load() should convert inline code to markdown")
	}
	if !strings.Contains(result, "```") {
		t.Errorf("Load() should convert code blocks to markdown")
	}
}

func TestHTMLConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewHTMLConverter()
	_, err := converter.Load("/nonexistent/file.html")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to read HTML file") {
		t.Errorf("Load() error should mention HTML file reading failure")
	}
}

func TestHTMLConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary HTML file with no read permissions
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "noperm.html")

	err := os.WriteFile(htmlFile, []byte("<p>test</p>"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	_, err = converter.Load(htmlFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestHTMLConverter_Load_InvalidHTML(t *testing.T) {
	// Create a temporary file with malformed HTML
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "invalid.html")

	// Malformed HTML should still be handled gracefully by the HTML-to-markdown converter
	htmlContent := `<p>Unclosed paragraph
<div>Unclosed div
<h1>Header without closing`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	// The html-to-markdown library should handle malformed HTML gracefully
	if err != nil {
		t.Errorf("Load() should handle malformed HTML gracefully, got error: %v", err)
	}

	// Result should contain some converted content
	if len(result) == 0 {
		t.Errorf("Load() should return some content even for malformed HTML")
	}
}

func TestHTMLConverter_Load_UnicodeHTML(t *testing.T) {
	// Create a temporary HTML file with Unicode characters
	tempDir := t.TempDir()
	htmlFile := filepath.Join(tempDir, "unicode.html")

	htmlContent := `<h1>こんにちは世界</h1>
<p>This is a test with emoji: 🚀🌟</p>
<p>Chinese: 你好世界</p>
<p>Arabic: مرحبا بالعالم</p>`

	err := os.WriteFile(htmlFile, []byte(htmlContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test HTML file: %v", err)
	}

	converter := NewHTMLConverter()
	result, err := converter.Load(htmlFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that Unicode characters are preserved
	if !strings.Contains(result, "こんにちは世界") {
		t.Errorf("Load() should preserve Japanese characters")
	}
	if !strings.Contains(result, "🚀🌟") {
		t.Errorf("Load() should preserve emoji")
	}
	if !strings.Contains(result, "你好世界") {
		t.Errorf("Load() should preserve Chinese characters")
	}
	if !strings.Contains(result, "مرحبا بالعالم") {
		t.Errorf("Load() should preserve Arabic characters")
	}
}
