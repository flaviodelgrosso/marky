package loaders

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCsvLoader_CanLoadMimeType(t *testing.T) {
	loader := &CsvLoader{}

	tests := []struct {
		name     string
		mimeType string
		expected bool
	}{
		{
			name:     "text/csv",
			mimeType: "text/csv",
			expected: true,
		},
		{
			name:     "application/csv",
			mimeType: "application/csv",
			expected: true,
		},
		{
			name:     "text/csv with charset",
			mimeType: "text/csv; charset=utf-8",
			expected: true,
		},
		{
			name:     "application/csv with parameters",
			mimeType: "application/csv; boundary=something",
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

func TestCsvLoader_Load_Success(t *testing.T) {
	tempDir := t.TempDir()

	tests := []struct {
		name       string
		csvContent string
		expectedMd string
	}{
		{
			name:       "Simple CSV",
			csvContent: "Name,Age,City\nJohn,30,New York\nJane,25,Los Angeles",
			expectedMd: "| Name | Age | City |\n| --- | --- | --- |\n| John | 30 | New York |\n| Jane | 25 | Los Angeles |\n",
		},
		{
			name: "CSV with quotes",
			csvContent: `Name,Description
John,"Software Engineer"
Jane,"Data Scientist"`,
			expectedMd: "| Name | Description |\n| --- | --- |\n| John | Software Engineer |\n| Jane | Data Scientist |\n",
		},
		{
			name: "CSV with commas in quotes",
			csvContent: `Name,Address
John,"123 Main St, New York"
Jane,"456 Oak Ave, Los Angeles"`,
			expectedMd: "| Name | Address |\n| --- | --- |\n| John | 123 Main St, New York |\n| Jane | 456 Oak Ave, Los Angeles |\n",
		},
		{
			name:       "Empty CSV",
			csvContent: "",
			expectedMd: "",
		},
		{
			name:       "Single column",
			csvContent: "Names\nJohn\nJane\nBob",
			expectedMd: "| Names |\n| --- |\n| John |\n| Jane |\n| Bob |\n",
		},
		{
			name:       "CSV with pipe characters",
			csvContent: "Name,Data\nJohn,value|with|pipes\nJane,another|pipe",
			expectedMd: "| Name | Data |\n| --- | --- |\n| John | value\\|with\\|pipes |\n| Jane | another\\|pipe |\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testFile := filepath.Join(tempDir, "test.csv")
			err := os.WriteFile(testFile, []byte(tt.csvContent), 0o644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			loader := &CsvLoader{}
			result, err := loader.Load(testFile)
			if err != nil {
				t.Errorf("Load() error = %v, wantErr false", err)
				return
			}

			if result != tt.expectedMd {
				t.Errorf("Load() = %v, want %v", result, tt.expectedMd)
			}
		})
	}
}

func TestCsvLoader_Load_FileNotFound(t *testing.T) {
	loader := &CsvLoader{}
	_, err := loader.Load("/nonexistent/file.csv")
	if err == nil {
		t.Error("Load() should return error for nonexistent file")
	}

	expectedError := "failed to load CSV file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedError)
	}
}

func TestCsvLoader_Load_InvalidCSV(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.csv")

	// Create a CSV with invalid format (unclosed quote)
	invalidCsv := `Name,Description
John,"Unclosed quote
Jane,Valid`

	err := os.WriteFile(testFile, []byte(invalidCsv), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	loader := &CsvLoader{}
	_, err = loader.Load(testFile)
	if err == nil {
		t.Error("Load() should return error for invalid CSV")
	}

	expectedError := "failed to load CSV file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Load() error = %v, should contain %v", err, expectedError)
	}
}

func TestReadCsvFile_Success(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.csv")
	csvContent := "Name,Age\nJohn,30\nJane,25"

	err := os.WriteFile(testFile, []byte(csvContent), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	records, err := readCsvFile(testFile)
	if err != nil {
		t.Errorf("readCsvFile() error = %v, wantErr false", err)
		return
	}

	expectedRecords := [][]string{
		{"Name", "Age"},
		{"John", "30"},
		{"Jane", "25"},
	}

	if len(records) != len(expectedRecords) {
		t.Errorf("readCsvFile() records length = %v, want %v", len(records), len(expectedRecords))
		return
	}

	for i, record := range records {
		if len(record) != len(expectedRecords[i]) {
			t.Errorf("readCsvFile() record %d length = %v, want %v", i, len(record), len(expectedRecords[i]))
			continue
		}
		for j, cell := range record {
			if cell != expectedRecords[i][j] {
				t.Errorf("readCsvFile() record[%d][%d] = %v, want %v", i, j, cell, expectedRecords[i][j])
			}
		}
	}
}

func TestReadCsvFile_EmptyFile(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "empty.csv")

	err := os.WriteFile(testFile, []byte(""), 0o644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	records, err := readCsvFile(testFile)
	if err != nil {
		t.Errorf("readCsvFile() error = %v, wantErr false", err)
		return
	}

	if len(records) != 0 {
		t.Errorf("readCsvFile() records length = %v, want 0", len(records))
	}
}

func TestReadCsvFile_FileNotFound(t *testing.T) {
	_, err := readCsvFile("/nonexistent/file.csv")
	if err == nil {
		t.Error("readCsvFile() should return error for nonexistent file")
	}

	expectedError := "unable to open file"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("readCsvFile() error = %v, should contain %v", err, expectedError)
	}
}
