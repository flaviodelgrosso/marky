package converters

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestNewExcelConverter(t *testing.T) {
	converter := NewExcelConverter()

	expectedExtensions := []string{".xlsx", ".xls"}
	expectedMimeTypes := []string{
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
		"application/vnd.openxmlformats-officedocument.spreadsheetml",
		"application/vnd.ms-excel",
	}

	if !reflect.DeepEqual(converter.AcceptedExtensions(), expectedExtensions) {
		t.Errorf("NewExcelConverter() extensions = %v, want %v", converter.AcceptedExtensions(), expectedExtensions)
	}

	if !reflect.DeepEqual(converter.AcceptedMimeTypes(), expectedMimeTypes) {
		t.Errorf("NewExcelConverter() mimeTypes = %v, want %v", converter.AcceptedMimeTypes(), expectedMimeTypes)
	}
}

func TestExcelConverter_Load_ValidFile(t *testing.T) {
	// Create a temporary Excel file
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "test.xlsx")

	// Create Excel file with test data
	f := excelize.NewFile()
	defer f.Close()

	// Set headers
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")
	f.SetCellValue("Sheet1", "C1", "City")

	// Set data rows
	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", 30)
	f.SetCellValue("Sheet1", "C2", "New York")

	f.SetCellValue("Sheet1", "A3", "Jane")
	f.SetCellValue("Sheet1", "B3", 25)
	f.SetCellValue("Sheet1", "C3", "Los Angeles")

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	expected := "| Name | Age | City |\n| --- | --- | --- |\n| John | 30 | New York |\n| Jane | 25 | Los Angeles |\n"
	if result != expected {
		t.Errorf("Load() = %v, want %v", result, expected)
	}
}

func TestExcelConverter_Load_EmptyFile(t *testing.T) {
	// Create a temporary empty Excel file
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "empty.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Empty Excel file should return empty markdown table
	if result != "" {
		t.Errorf("Load() with empty Excel file = %v, want empty string", result)
	}
}

func TestExcelConverter_Load_OnlyHeader(t *testing.T) {
	// Create a temporary Excel file with only header
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "header_only.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	// Set only headers
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")
	f.SetCellValue("Sheet1", "C1", "City")

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	expected := "| Name | Age | City |\n| --- | --- | --- |\n"
	if result != expected {
		t.Errorf("Load() = %v, want %v", result, expected)
	}
}

func TestExcelConverter_Load_WithFormulas(t *testing.T) {
	// Create a temporary Excel file with formulas
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "formulas.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	// Set headers
	f.SetCellValue("Sheet1", "A1", "Value1")
	f.SetCellValue("Sheet1", "B1", "Value2")
	f.SetCellValue("Sheet1", "C1", "Sum")

	// Set data with formula
	f.SetCellValue("Sheet1", "A2", 10)
	f.SetCellValue("Sheet1", "B2", 20)
	f.SetCellFormula("Sheet1", "C2", "=A2+B2")

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that the result contains the calculated value or formula
	if !strings.Contains(result, "Value1") || !strings.Contains(result, "10") {
		t.Errorf("Load() should contain Excel data")
	}
}

func TestExcelConverter_Load_WithSpecialCharacters(t *testing.T) {
	// Create a temporary Excel file with special characters
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "special.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	// Set data with special characters including pipes
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Description")

	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", "Works at Company|Inc")

	f.SetCellValue("Sheet1", "A3", "Jane")
	f.SetCellValue("Sheet1", "B3", "Has pipe | character")

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
	if err != nil {
		t.Errorf("Load() returned unexpected error: %v", err)
	}

	// Check that pipe characters are properly escaped in markdown
	if !strings.Contains(result, `Company\|Inc`) {
		t.Errorf("Load() should escape pipe characters")
	}
	if !strings.Contains(result, `pipe \| character`) {
		t.Errorf("Load() should escape pipe characters")
	}
}

func TestExcelConverter_Load_UnicodeCharacters(t *testing.T) {
	// Create a temporary Excel file with Unicode characters
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "unicode.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	// Set data with Unicode characters
	f.SetCellValue("Sheet1", "A1", "名前")
	f.SetCellValue("Sheet1", "B1", "年齢")
	f.SetCellValue("Sheet1", "C1", "住所")

	f.SetCellValue("Sheet1", "A2", "佐藤太郎")
	f.SetCellValue("Sheet1", "B2", 30)
	f.SetCellValue("Sheet1", "C2", "東京")

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	converter := NewExcelConverter()
	result, err := converter.Load(excelFile)
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

func TestExcelConverter_Load_NonExistentFile(t *testing.T) {
	converter := NewExcelConverter()
	_, err := converter.Load("/nonexistent/file.xlsx")

	if err == nil {
		t.Errorf("Load() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "failed to load Excel file") {
		t.Errorf("Load() error should mention Excel file loading failure")
	}
}

func TestExcelConverter_Load_InvalidFile(t *testing.T) {
	// Create a temporary file that's not a valid Excel file
	tempDir := t.TempDir()
	invalidFile := filepath.Join(tempDir, "invalid.xlsx")

	err := os.WriteFile(invalidFile, []byte("not an excel file"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	converter := NewExcelConverter()
	_, err = converter.Load(invalidFile)

	if err == nil {
		t.Errorf("Load() should return error for invalid Excel file")
	}

	if !strings.Contains(err.Error(), "failed to load Excel file") {
		t.Errorf("Load() error should mention Excel file loading failure")
	}
}

func TestReadExcelFile_ValidFile(t *testing.T) {
	// Create a temporary Excel file
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "test.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	// Set test data
	f.SetCellValue("Sheet1", "A1", "Name")
	f.SetCellValue("Sheet1", "B1", "Age")
	f.SetCellValue("Sheet1", "A2", "John")
	f.SetCellValue("Sheet1", "B2", 30)

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	rows, err := readExcelFile(excelFile)
	if err != nil {
		t.Errorf("readExcelFile() returned unexpected error: %v", err)
	}

	expectedRows := [][]string{
		{"Name", "Age"},
		{"John", "30"},
	}

	if !reflect.DeepEqual(rows, expectedRows) {
		t.Errorf("readExcelFile() = %v, want %v", rows, expectedRows)
	}
}

func TestReadExcelFile_NonExistentFile(t *testing.T) {
	_, err := readExcelFile("/nonexistent/file.xlsx")

	if err == nil {
		t.Errorf("readExcelFile() should return error for non-existent file")
	}

	if !strings.Contains(err.Error(), "unable to open Excel file") {
		t.Errorf("readExcelFile() error should mention Excel file opening failure")
	}
}

func TestReadExcelFile_EmptyWorkbook(t *testing.T) {
	// Create an Excel file with no sheets (though this is unlikely in practice)
	// We'll test with a regular empty sheet instead
	tempDir := t.TempDir()
	excelFile := filepath.Join(tempDir, "empty.xlsx")

	f := excelize.NewFile()
	defer f.Close()

	err := f.SaveAs(excelFile)
	if err != nil {
		t.Fatalf("Failed to create test Excel file: %v", err)
	}

	rows, err := readExcelFile(excelFile)
	if err != nil {
		t.Errorf("readExcelFile() returned unexpected error: %v", err)
	}

	// Empty sheet should return empty rows
	if len(rows) != 0 {
		t.Errorf("readExcelFile() with empty sheet should return empty rows, got %d rows", len(rows))
	}
}
