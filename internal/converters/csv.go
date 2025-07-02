package converters

import (
	"encoding/csv"
	"fmt"
	"os"

	"github.com/flaviodelgrosso/marky/internal/utils"
)

// CsvConverter handles loading and converting CSV files to markdown tables.
type CsvConverter struct {
	BaseConverter
}

// NewCsvConverter creates a new CSV converter with appropriate MIME types and extensions.
func NewCsvConverter() Converter {
	return &CsvConverter{
		BaseConverter: NewBaseConverter(
			[]string{".csv"},
			[]string{"text/csv", "application/csv"},
		),
	}
}

// Load reads a CSV file and converts it to a markdown table.
func (*CsvConverter) Load(path string) (string, error) {
	records, err := readCsvFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load CSV file: %w", err)
	}

	return utils.ToMarkdownTable(records), nil
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
