package loaders

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEpubLoader_CanLoadMimeType(t *testing.T) {
	loader := &EpubLoader{}

	testCases := []struct {
		mimeType string
		expected bool
	}{
		{"application/epub+zip", true},
		{"application/epub", true},
		{"application/x-epub+zip", true},
		{"application/pdf", false},
		{"text/html", false},
		{"application/vnd.openxmlformats-officedocument.wordprocessingml.document", false},
	}

	for _, tc := range testCases {
		t.Run(tc.mimeType, func(t *testing.T) {
			result := loader.CanLoadMimeType(tc.mimeType)
			if result != tc.expected {
				t.Errorf("CanLoadMimeType(%q) = %v, expected %v", tc.mimeType, result, tc.expected)
			}
		})
	}
}

func TestEpubLoader_Load(t *testing.T) {
	// Create a temporary EPUB file for testing
	tempDir := t.TempDir()
	epubPath := filepath.Join(tempDir, "test.epub")

	if err := createTestEpub(epubPath); err != nil {
		t.Fatalf("Failed to create test EPUB: %v", err)
	}

	loader := &EpubLoader{}
	result, err := loader.Load(epubPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if result == "" {
		t.Error("Load() returned empty result")
	}

	// Check if the result contains expected metadata
	if !strings.Contains(result, "**Title:**") {
		t.Error("Result should contain title metadata")
	}

	// Check if the result contains expected content
	if !strings.Contains(result, "Test Chapter") {
		t.Error("Result should contain chapter content")
	}
}

func TestEpubLoader_Load_InvalidFile(t *testing.T) {
	loader := &EpubLoader{}
	_, err := loader.Load("nonexistent.epub")

	if err == nil {
		t.Error("Load() should fail for nonexistent file")
	}
}

// createTestEpub creates a minimal valid EPUB file for testing
func createTestEpub(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	// Create mimetype file
	mimetypeFile, err := w.Create("mimetype")
	if err != nil {
		return err
	}
	if _, err := mimetypeFile.Write([]byte("application/epub+zip")); err != nil {
		return err
	}

	// Create META-INF/container.xml
	containerFile, err := w.Create("META-INF/container.xml")
	if err != nil {
		return err
	}
	containerXML := `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
  <rootfiles>
    <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
  </rootfiles>
</container>`
	if _, err := containerFile.Write([]byte(containerXML)); err != nil {
		return err
	}

	// Create OEBPS/content.opf
	opfFile, err := w.Create("OEBPS/content.opf")
	if err != nil {
		return err
	}
	opfXML := `<?xml version="1.0" encoding="UTF-8"?>
<package xmlns="http://www.idpf.org/2007/opf" version="3.0" unique-identifier="bookid">
  <metadata xmlns:dc="http://purl.org/dc/elements/1.1/">
    <dc:title>Test Book</dc:title>
    <dc:creator>Test Author</dc:creator>
    <dc:language>en</dc:language>
    <dc:identifier id="bookid">test-book-123</dc:identifier>
  </metadata>
  <manifest>
    <item id="chapter1" href="chapter1.xhtml" media-type="application/xhtml+xml"/>
  </manifest>
  <spine>
    <itemref idref="chapter1"/>
  </spine>
</package>`
	if _, err := opfFile.Write([]byte(opfXML)); err != nil {
		return err
	}

	// Create OEBPS/chapter1.xhtml
	chapterFile, err := w.Create("OEBPS/chapter1.xhtml")
	if err != nil {
		return err
	}
	chapterXHTML := `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html>
<html xmlns="http://www.w3.org/1999/xhtml">
<head>
  <title>Chapter 1</title>
</head>
<body>
  <h1>Test Chapter</h1>
  <p>This is a test paragraph in the EPUB file.</p>
  <p>Another paragraph with <strong>bold text</strong> and <em>italic text</em>.</p>
</body>
</html>`
	if _, err := chapterFile.Write([]byte(chapterXHTML)); err != nil {
		return err
	}

	return nil
}
