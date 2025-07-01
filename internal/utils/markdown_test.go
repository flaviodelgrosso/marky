package utils

import (
	"strings"
	"testing"
)

func TestToMarkdownTable_EmptyInput(t *testing.T) {
	result := ToMarkdownTable([][]string{})
	if result != "" {
		t.Errorf("ToMarkdownTable() with empty input = %v, want empty string", result)
	}
}

func TestToMarkdownTable_SingleRow(t *testing.T) {
	input := [][]string{
		{"Name", "Age", "City"},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Age | City |\n| --- | --- | --- |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_MultipleRows(t *testing.T) {
	input := [][]string{
		{"Name", "Age", "City"},
		{"John", "30", "New York"},
		{"Jane", "25", "Los Angeles"},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Age | City |\n| --- | --- | --- |\n| John | 30 | New York |\n| Jane | 25 | Los Angeles |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_EscapePipeCharacters(t *testing.T) {
	input := [][]string{
		{"Name", "Description"},
		{"John", "Works at Company|Inc"},
		{"Jane", "Has pipe | character"},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Description |\n| --- | --- |\n| John | Works at Company\\|Inc |\n| Jane | Has pipe \\| character |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_TrimWhitespace(t *testing.T) {
	input := [][]string{
		{"  Name  ", " Age ", "City   "},
		{" John ", "30  ", "  New York "},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Age | City |\n| --- | --- | --- |\n| John | 30 | New York |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_UnevenRows(t *testing.T) {
	input := [][]string{
		{"Name", "Age", "City", "Country"},
		{"John", "30", "New York"},               // Missing country
		{"Jane", "25"},                           // Missing city and country
		{"Bob", "35", "Chicago", "USA", "Extra"}, // Extra column
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Age | City | Country |\n| --- | --- | --- | --- |\n| John | 30 | New York |  |\n| Jane | 25 |  |  |\n| Bob | 35 | Chicago | USA |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_EmptyStrings(t *testing.T) {
	input := [][]string{
		{"Name", "Age", "City"},
		{"John", "", "New York"},
		{"", "25", ""},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Age | City |\n| --- | --- | --- |\n| John |  | New York |\n|  | 25 |  |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_SpecialCharacters(t *testing.T) {
	input := [][]string{
		{"Name", "Special Characters"},
		{"John", "Hello & goodbye"},
		{"Jane", "< > \" ' & *"},
		{"Bob", "Line\nbreak"},
	}

	result := ToMarkdownTable(input)

	// Should escape pipes but preserve other characters
	if !strings.Contains(result, "Hello & goodbye") {
		t.Error("ToMarkdownTable() should preserve ampersand characters")
	}

	if !strings.Contains(result, "< > \" ' & *") {
		t.Error("ToMarkdownTable() should preserve special characters except pipes")
	}

	if !strings.Contains(result, "Line\nbreak") {
		t.Error("ToMarkdownTable() should preserve newline characters")
	}
}

func TestToMarkdownTable_LongContent(t *testing.T) {
	longText := strings.Repeat("This is a very long text content that should be handled properly. ", 10)
	input := [][]string{
		{"ID", "Content"},
		{"1", longText},
	}

	result := ToMarkdownTable(input)

	// Should contain the long text
	if !strings.Contains(result, longText) {
		t.Error("ToMarkdownTable() should handle long content properly")
	}

	// Should still have proper table structure
	lines := strings.Split(result, "\n")
	if len(lines) < 4 { // header, separator, data row, empty line
		t.Error("ToMarkdownTable() should maintain table structure with long content")
	}
}

func TestToMarkdownTable_UnicodeCharacters(t *testing.T) {
	input := [][]string{
		{"Name", "Unicode"},
		{"🚀", "Rocket"},
		{"数据", "Chinese"},
		{"café", "French"},
	}

	result := ToMarkdownTable(input)
	expected := "| Name | Unicode |\n| --- | --- |\n| 🚀 | Rocket |\n| 数据 | Chinese |\n| café | French |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() = %v, want %v", result, expected)
	}
}

func TestToMarkdownTable_OnlyHeader(t *testing.T) {
	input := [][]string{
		{"Column1", "Column2", "Column3"},
	}

	result := ToMarkdownTable(input)
	expected := "| Column1 | Column2 | Column3 |\n| --- | --- | --- |\n"

	if result != expected {
		t.Errorf("ToMarkdownTable() with only header = %v, want %v", result, expected)
	}
}
