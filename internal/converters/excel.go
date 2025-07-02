package converters

import (
	"fmt"
	"log"

	"github.com/flaviodelgrosso/marky/internal/utils"
	"github.com/xuri/excelize/v2"
)

// ExcelConverter handles loading and converting Excel files to markdown tables.
type ExcelConverter struct {
	BaseConverter
}

// NewExcelConverter creates a new Excel converter with appropriate MIME types and extensions.
func NewExcelConverter() Converter {
	return &ExcelConverter{
		BaseConverter: NewBaseConverter(
			[]string{".xlsx", ".xls"},
			[]string{
				"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				"application/vnd.openxmlformats-officedocument.spreadsheetml",
				"application/vnd.ms-excel",
			},
		),
	}
}

// Load reads an Excel file and converts it to a markdown table.
func (*ExcelConverter) Load(path string) (string, error) {
	rows, err := readExcelFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to load Excel file: %w", err)
	}

	return utils.ToMarkdownTable(rows), nil
}

// readExcelFile reads and parses an Excel file, returning all records from the first sheet.
func readExcelFile(path string) ([][]string, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to open Excel file %s: %w", path, err)
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			// Log the close error, but don't override the main error
			log.Printf("Warning: failed to close Excel file %s: %v\n", path, closeErr)
		}
	}()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file %s", path)
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("unable to read rows from sheet %s in file %s: %w", sheets[0], path, err)
	}

	return rows, nil
}
