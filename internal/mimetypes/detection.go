package mimetypes

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// MimeTypeInfo contains MIME type and file extension information
type MimeTypeInfo struct {
	MimeType  string
	Extension string
}

// DetectMimeType detects the MIME type of a file from its content and extension
func DetectMimeType(filePath string) (*MimeTypeInfo, error) {
	// Read the file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file for MIME type detection: %w", err)
	}

	// Use only first 512 bytes for MIME type detection
	sampleSize := min(len(data), 512)
	sample := data[:sampleSize]

	// Get the file extension
	ext := strings.ToLower(filepath.Ext(filePath))

	// Detect MIME type from content
	contentType := http.DetectContentType(sample)

	// Override with more specific detection for Office documents and other formats
	mimeType := enhanceMimeTypeDetection(sample, ext, contentType)

	return &MimeTypeInfo{
		MimeType:  mimeType,
		Extension: ext,
	}, nil
}

// Magic numbers for file type detection
var (
	zipSignature = []byte{0x50, 0x4B, 0x03, 0x04}                         // PK..
	pdfSignature = []byte{0x25, 0x50, 0x44, 0x46}                         // %PDF
	oleSignature = []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1} // OLE Compound File
)

// MIME type constants
const (
	mimeDocx  = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	mimeXlsx  = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	mimePptx  = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	mimePdf   = "application/pdf"
	mimeDoc   = "application/msword"
	mimeXls   = "application/vnd.ms-excel"
	mimePpt   = "application/vnd.ms-powerpoint"
	mimeCsv   = "text/csv"
	mimeHTML  = "text/html"
	mimeXML   = "application/xml"
	mimeIpynb = "application/x-ipynb+json"
	mimeEpub  = "application/epub+zip"
)

// enhanceMimeTypeDetection provides better detection for specific file types
func enhanceMimeTypeDetection(data []byte, ext, detectedType string) string {
	// Check for signatures of known file types
	if bytes.HasPrefix(data, zipSignature) { // Office Open XML formats (DOCX, XLSX, PPTX) and EPUB are ZIP based
		switch ext {
		case ".docx":
			return mimeDocx
		case ".xlsx":
			return mimeXlsx
		case ".pptx":
			return mimePptx
		case ".epub":
			return mimeEpub
		}
	}

	if bytes.HasPrefix(data, pdfSignature) {
		return mimePdf
	}

	if bytes.HasPrefix(data, oleSignature) { // Older Office formats
		switch ext {
		case ".doc":
			return mimeDoc
		case ".xls":
			return mimeXls
		case ".ppt":
			return mimePpt
		}
	}

	// For text-based formats, trust the extension if the content seems to be text.
	if strings.HasPrefix(detectedType, "text/") {
		switch ext {
		case ".csv":
			return mimeCsv
		case ".html", ".htm":
			return mimeHTML
		case ".xml":
			return mimeXML
		case ".ipynb":
			return mimeIpynb
		}
	}

	return detectedType
}

// IsMimeTypeSupported checks if a MIME type matches any of the given patterns
func IsMimeTypeSupported(mimeType string, supportedTypes []string) bool {
	for _, supportedType := range supportedTypes {
		if strings.HasPrefix(mimeType, supportedType) {
			return true
		}
	}
	return false
}
