package utils

import "testing"

func TestRuneWidth(t *testing.T) {
	tests := []struct {
		rune     rune
		expected int
		name     string
	}{
		{rune(0), 0, "NUL control char"},
		{rune(31), 0, "Unit separator control char"},
		{rune(127), 0, "DEL control char"},
		{'A', 1, "ASCII letter"},
		{' ', 1, "ASCII space"},
		{'~', 1, "ASCII tilde"},
		{rune(0x200B), 0, "Zero Width Space"},
		{rune(0x200C), 0, "Zero Width Non-Joiner"},
		{rune(0x200D), 0, "Zero Width Joiner"},
		{rune(0xFEFF), 0, "Zero Width No-Break Space (BOM)"},
		{rune(0x0301), 0, "Combining acute accent"},
		{rune(0x1F600), 2, "Emoji (grinning face)"},
		{rune(0x4E2D), 2, "CJK Unified Ideograph (中)"},
		{rune(0xFF21), 2, "Fullwidth Latin Capital Letter A"},
		{rune(0xFF66), 1, "Halfwidth Katakana Wo"},
		{rune(0xAC00), 2, "Hangul Syllable"},
		{rune(0x3042), 2, "Hiragana (あ)"},
		{rune(0x30A2), 2, "Katakana (ア)"},
		{rune(0xFF9E), 1, "Halfwidth Katakana Voiced Sound Mark"},
	}

	for _, tt := range tests {
		got := RuneWidth(tt.rune)
		if got != tt.expected {
			t.Errorf("RuneWidth(%U) [%s] = %d; want %d", tt.rune, tt.name, got, tt.expected)
		}
	}
}

func TestStringWidth(t *testing.T) {
	cases := []struct {
		input    string
		expected int
		name     string
	}{
		{"Hello", 5, "ASCII only"},
		{"A中B", 4, "Mixed ASCII and CJK"},
		{"A\u200B中B", 4, "Zero width in string"},
		{"A\u0301B", 2, "Combining mark in string"},
		{"\u1F600\u1F601", 4, "Two emojis"},
		{"\uFF21\uFF66", 3, "Fullwidth and halfwidth"},
	}

	for _, c := range cases {
		got := StringWidth(c.input)
		if got != c.expected {
			t.Errorf("StringWidth(%q) [%s] = %d; want %d", c.input, c.name, got, c.expected)
		}
	}
}
