package utils

import (
	"unicode"
)

// RuneWidth returns the display width of a single rune
// Returns:
//   - 0 for control characters, combining marks, and zero-width characters
//   - 1 for normal ASCII and most Unicode characters
//   - 2 for wide characters (CJK ideographs, fullwidth forms, etc.)
func RuneWidth(r rune) int {
	// Handle basic ASCII control characters
	if r < 32 || r == 127 {
		return 0
	}

	// Handle basic ASCII printable characters
	if r < 127 {
		return 1
	}

	// Handle common zero-width characters
	switch r {
	case 0x200B, // Zero Width Space
		0x200C, // Zero Width Non-Joiner
		0x200D, // Zero Width Joiner
		0xFEFF: // Zero Width No-Break Space (BOM)
		return 0
	}

	// Handle combining marks (they don't add width)
	if unicode.In(r, unicode.Mn, unicode.Me, unicode.Mc) {
		return 0
	}

	// Check for wide characters
	if isWideRune(r) {
		return 2
	}

	// Default to width 1 for other printable characters
	return 1
}

// wideRanges defines Unicode ranges for wide characters (2-column display width)
var wideRanges = [][2]rune{
	{0x1F300, 0x1F5FF}, // Miscellaneous Symbols and Pictographs
	{0x1F600, 0x1F64F}, // Emoticons
	{0x1F680, 0x1F6FF}, // Transport and Map Symbols
	{0x1F700, 0x1F77F}, // Alchemical Symbols
	{0x1F780, 0x1F7FF}, // Geometric Shapes Extended
	{0x1F800, 0x1F8FF}, // Supplemental Arrows-C
	{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
	{0x20000, 0x2A6DF}, // CJK Extension B and beyond
	{0x3000, 0x303F},   // CJK Symbols and Punctuation
	{0x3040, 0x309F},   // Hiragana
	{0x30A0, 0x30FF},   // Katakana
	{0x3400, 0x4DBF},   // CJK Extension A
	{0x4E00, 0x9FFF},   // CJK Unified Ideographs
	{0xAC00, 0xD7AF},   // Hangul Syllables
	{0xFF01, 0xFF60},   // Fullwidth ASCII variants
	{0xFFE0, 0xFFE6},   // Fullwidth symbols
}

// halfWidthRanges defines Unicode ranges that are explicitly half-width
var halfWidthRanges = [][2]rune{
	{0xFF61, 0xFFDC}, // Halfwidth and Fullwidth Forms (halfwidth part)
}

// isWideRune determines if a rune should be considered "wide" (2 columns)
func isWideRune(r rune) bool {
	// Check half-width ranges first (explicit narrow characters)
	for _, rang := range halfWidthRanges {
		if r >= rang[0] && r <= rang[1] {
			return false
		}
	}

	// Check wide ranges
	for _, rang := range wideRanges {
		if r >= rang[0] && r <= rang[1] {
			return true
		}
	}

	return false
}

// StringWidth returns the display width of a string
func StringWidth(s string) int {
	width := 0
	for _, r := range s {
		width += RuneWidth(r)
	}
	return width
}
