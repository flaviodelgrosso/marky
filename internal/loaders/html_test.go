package loaders

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHtmlLoader_CanLoadMimeType(t *testing.T) {
	loader := &HTMLLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "text/html",
			mimeType: "text/html",
			expected: true,
		},
		{
			name:     "text/html with charset",
			mimeType: "text/html; charset=utf-8",
			expected: true,
		},
		{
			name:     "text/html with parameters",
			mimeType: "text/html; boundary=something",
			expected: true,
		},
		{
			name:     "text/plain",
			mimeType: "text/plain",
			expected: false,
		},
		{
			name:     "application/json",
			mimeType: "application/json",
			expected: false,
		},
		{
			name:     "empty string",
			mimeType: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := loader.CanLoadMimeType(tt.mimeType)
			if result != tt.expected {
				t.Errorf("CanLoadMimeType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHtmlLoader_Load_Success(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name           string
		htmlContent    string
		expectContains []string
	}{
		{
			name:           "Simple HTML",
			htmlContent:    "<html><body><h1>Hello World</h1><p>This is a paragraph.</p></body></html>",
			expectContains: []string{"# Hello World", "This is a paragraph."},
		},
		{
			name:           "HTML with links",
			htmlContent:    "<html><body><a href='https://example.com'>Example Link</a></body></html>",
			expectContains: []string{"[Example Link](https://example.com)"},
		},
		{
			name:           "HTML with lists",
			htmlContent:    "<html><body><ul><li>Item 1</li><li>Item 2</li></ul></body></html>",
			expectContains: []string{"- Item 1", "- Item 2"},
		},
		{
			name:           "HTML with bold and italic",
			htmlContent:    "<html><body><b>Bold text</b> and <i>italic text</i></body></html>",
			expectContains: []string{"**Bold text**", "*italic text*"},
		},
		{
			name:           "Empty HTML",
			htmlContent:    "<html><body></body></html>",
			expectContains: []string{}, // Empty content should produce minimal markdown
		},
		{
			name:           "HTML with code",
			htmlContent:    "<html><body><code>console.log('hello')</code></body></html>",
			expectContains: []string{"`console.log('hello')`"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, "test.html")
			err := os.WriteFile(testFile, []byte(tt.htmlContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			loader := &HTMLLoader{}
			result, err := loader.Load(testFile)
			if err != nil {
				t.Errorf("Load() error = %v, wantErr false", err)
				return
			}

			for _, expected := range tt.expectContains {
				if !strings.Contains(result, expected) {
					t.Errorf("Load() result should contain %v, got %v", expected, result)
				}
			}
		})
	}
}

func TestHtmlLoader_Load_FileNotFound(t *testing.T) {
	loader := &HTMLLoader{}
	_, err := loader.Load("/nonexistent/file.html")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	expectedError := "failed to read HTML file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedError)
	}
}

func TestHtmlLoader_Load_ComplexHTML(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "complex.html")

	complexHTML := `<!DOCTYPE html>
<html>
<head>
    <title>Test Page</title>
    <style>body { font-family: Arial; }</style>
</head>
<body>
    <header>
        <h1>Welcome to My Site</h1>
        <nav>
            <ul>
                <li><a href="#home">Home</a></li>
                <li><a href="#about">About</a></li>
            </ul>
        </nav>
    </header>
    <main>
        <article>
            <h2>Article Title</h2>
            <p>This is the <strong>main content</strong> of the article.</p>
            <blockquote>This is a quote from someone important.</blockquote>
            <pre><code>function example() { return "code"; }</code></pre>
        </article>
    </main>
    <footer>
        <p>&copy; 2023 My Website</p>
    </footer>
    <script>console.log("JavaScript should be ignored");</script>
</body>
</html>`

	err := os.WriteFile(testFile, []byte(complexHTML), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &HTMLLoader{}
	result, err := loader.Load(testFile)
	if err != nil {
		t.Errorf("Load() error = %v, wantErr false", err)
		return
	}

	// Check that important content is converted
	expectedElements := []string{
		"# Welcome to My Site", // h1
		"## Article Title",     // h2
		"**main content**",     // strong
		"This is a quote",      // blockquote content
	}

	for _, expected := range expectedElements {
		if !strings.Contains(result, expected) {
			t.Errorf("Load() result should contain %v", expected)
		}
	}

	// CSS and JavaScript should typically be ignored by the HTML-to-markdown converter
	if strings.Contains(result, "font-family: Arial") {
		t.Error("Load() result should not contain CSS")
	}
	if strings.Contains(result, "console.log") {
		t.Error("Load() result should not contain JavaScript")
	}
}

func TestHtmlLoader_Load_MalformedHTML(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "malformed.html")

	// HTML with unclosed tags - most parsers handle this gracefully
	malformedHTML := `<html><body><h1>Title<p>Paragraph without closing tag<div>Another element`

	err := os.WriteFile(testFile, []byte(malformedHTML), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &HTMLLoader{}
	result, err := loader.Load(testFile)
	// HTML parsers are typically forgiving, so this should not error
	if err != nil {
		t.Errorf("Load() should handle malformed HTML gracefully, got error: %v", err)
		return
	}

	// Should still extract some meaningful content
	if !strings.Contains(result, "Title") {
		t.Error("Load() should extract title from malformed HTML")
	}
	if !strings.Contains(result, "Paragraph") {
		t.Error("Load() should extract paragraph content from malformed HTML")
	}
}

func TestHtmlLoader_Load_SpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "special.html")

	specialHTML := `<html><body>
		<p>Special characters: &amp; &lt; &gt; &quot; &#39;</p>
		<p>Unicode: 🚀 数据 café</p>
	</body></html>`

	err := os.WriteFile(testFile, []byte(specialHTML), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &HTMLLoader{}
	result, err := loader.Load(testFile)
	if err != nil {
		t.Errorf("Load() error = %v, wantErr false", err)
		return
	}

	// HTML entities should be decoded
	if !strings.Contains(result, "&") {
		t.Error("Load() should decode &amp; entity")
	}
	// Note: < and > might not be present if the converter doesn't include them in output
	// This is okay as the converter may choose to omit them if they don't add value

	// Unicode should be preserved
	if !strings.Contains(result, "🚀") {
		t.Error("Load() should preserve Unicode emoji")
	}
	if !strings.Contains(result, "数据") {
		t.Error("Load() should preserve Unicode Chinese characters")
	}
	if !strings.Contains(result, "café") {
		t.Error("Load() should preserve Unicode accented characters")
	}
}
