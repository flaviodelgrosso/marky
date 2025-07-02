package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestNewCsvConverter(t *testing.T) {
	converter := NewCsvConverter()

	expectedExtensions := []string{".csv"}
	expectedMimeTypes := []string{"text/csv", "application/csv"}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewCsvConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewCsvConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestCsvConverter_Load_ValidFile(t *testing.T) {
	// Create a temporary CSV file
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")

	csvContent := `Name,Age,City
John,30,New York
Jane,25,Los Angeles`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	expected := "| Name | Age | City |\n| --- | --- | --- |\n| John | 30 | New York |\n| Jane | 25 | Los Angeles |\n"
	if result != expected {
		t.Errorf("Load() = %v, want %v", result, expected)
	}
}

func TestCsvConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty CSV file
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "empty.csv")

	err := os.WriteFile(csvFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Empty CSV should return empty markdown table
	if result != "" {
		t.Errorf("Load() with empty CSV = %v, want empty string", result)
	}
}

func TestCsvConverter_Load_OnlyHeader(t *testing.T) {
	// Create a temporary CSV file with only header
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "header_only.csv")

	csvContent := "Name,Age,City"

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	expected := "| Name | Age | City |\n| --- | --- | --- |\n"
	if result != expected {
		t.Errorf("Load() = %v, want %v", result, expected)
	}
}

func TestCsvConverter_Load_WithQuotes(t *testing.T) {
	// Create a temporary CSV file with quoted fields
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "quoted.csv")

	csvContent := `"Name","Age","Description"
"John Doe",30,"Works at ""Company Inc"""
"Jane Smith",25,"Has, comma in description"`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that quoted content is properly handled
	if !strings.Contains(result, "John Doe") {
		t.Errorf("Load() should contain 'John Doe'")
	}
	if !strings.Contains(result, `Works at "Company Inc"`) {
		t.Errorf("Load() should contain properly unquoted description")
	}
	if !strings.Contains(result, "Has, comma in description") {
		t.Errorf("Load() should contain description with comma")
	}
}

func TestCsvConverter_Load_WithPipeCharacters(t *testing.T) {
	// Create a temporary CSV file with pipe characters
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "pipes.csv")

	csvContent := `Name,Description
John,Works at Company|Inc
Jane,Has pipe | character`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that pipe characters are properly escaped
	if !strings.Contains(result, `Company\|Inc`) {
		t.Errorf("Load() should escape pipe characters")
	}
	if !strings.Contains(result, `pipe \| character`) {
		t.Errorf("Load() should escape pipe characters")
	}
}

func TestCsvConverter_Load_UnevenRows(t *testing.T) {
	// The standard CSV parser requires all rows to have the same number of fields
	// This test verifies that uneven CSV files return an appropriate error
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "uneven.csv")

	csvContent := `Name,Age,City,Country
John,30,New York
Jane,25
Bob,35,Chicago,USA,Extra`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	_, err = converter.Load(csvFile)

	// The CSV parser should return an error for uneven rows
	if err == nil {
		t.Errorf("Load() should return error for uneven CSV rows")
	}

	if !strings.Contains(err.Error(), "failed to load CSV file") {
		t.Errorf("Load() error should mention CSV loading failure, got: %v", err)
	}
}

func TestCsvConverter_Load_UnicodeCharacters(t *testing.T) {
	// Create a temporary CSV file with Unicode characters
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "unicode.csv")

	csvContent := `名前,年齢,住所
佐藤太郎,30,東京
三木英子,25,大阪`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	result, err := converter.Load(csvFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that Unicode characters are preserved
	if !strings.Contains(result, "佐藤太郎") {
		t.Errorf("Load() should preserve Unicode characters")
	}
	if !strings.Contains(result, "東京") {
		t.Errorf("Load() should preserve Unicode characters")
	}
}

func TestCsvConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewCsvConverter()
	_, err := converter.Load("/nonexistent/file.csv")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to load CSV file") {
		t.Errorf("Load() error should mention CSV file loading failure")
	}
}

func TestCsvConverter_Load_InvalidPermissions(t *testing.T) {
	// Create a temporary CSV file with no read permissions
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "noperm.csv")

	err := os.WriteFile(csvFile, []byte("test,data"), 0o000)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	converter := NewCsvConverter()
	_, err = converter.Load(csvFile)

	if err == nil {
		t.Errorf("Load() should return error for file with no read permissions")
	}
}

func TestReadCsvFile_ValidFile(t *testing.T) {
	// Create a temporary CSV file
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "test.csv")

	csvContent := `Name,Age,City
John,30,New York
Jane,25,Los Angeles`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	records, err := readCsvFile(csvFile)
	if err != nil {
		t.Errorf("readCsvFile() returned unexpected error: %v", err)
	}

	expectedRecords := [][]string{
		{"Name", "Age", "City"},
		{"John", "30", "New York"},
		{"Jane", "25", "Los Angeles"},
	}

	if !reflect.DeepEqual(records, expectedRecords) {
		t.Errorf("readCsvFile() = %v, want %v", records, expectedRecords)
	}
}

func TestReadCsvFile_NonExistentFile(t *testing.T) {
	_, err := readCsvFile("/nonexistent/file.csv")

	if err == nil {
		t.Errorf("readCsvFile() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "unable to open file") {
		t.Errorf("readCsvFile() error should mention file opening failure")
	}
}

func TestReadCsvFile_InvalidCSV(t *testing.T) {
	// Create a temporary file with invalid CSV content
	tempDir := t.TempDir()
	csvFile := filepath.Join(tempDir, "invalid.csv")

	// CSV with unclosed quotes
	csvContent := `"Name,"Age","City"
"John","30","New York`

	err := os.WriteFile(csvFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	_, err = readCsvFile(csvFile)

	if err == nil {
		t.Errorf("readCsvFile() should return error for invalid CSV")
	}

	if !strings.Contains(err.Error(), "unable to parse CSV file") {
		t.Errorf("readCsvFile() error should mention CSV parsing failure")
	}
}
