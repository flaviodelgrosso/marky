package loaders

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/flaviodelgrosso/marky/internal/mimetypes"
	"github.com/flaviodelgrosso/marky/internal/utils"
)

// CsvLoader handles loading and converting CSV files to markdown tables.
type CsvLoader struct{}

// Load reads a CSV file and converts it to a markdown table.
func (*CsvLoader) Load(path string) (string, error) {
	records, err := readCsvFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load CSV file: %w", err)
	}

	return utils.ToMarkdownTable(records), nil
}

// CanLoadMimeType returns true if the MIME type is supported for CSV files.
func (*CsvLoader) CanLoadMimeType(mimeType string) bool {
	supportedTypes := []string{
		"text/csv",
		"application/csv",
	}
	return mimetypes.IsMimeTypeSupported(mimeType, supportedTypes)
}

// readCsvFile reads and parses a CSV file, returning all records.
func readCsvFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open file %s: %w", path, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("unable to parse CSV file %s: %w", path, err)
	}

	return records, nil
}
