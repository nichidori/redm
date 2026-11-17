package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func FixLength(v string, maxLength int) string {
	if utf8.RuneCountInString(v) <= maxLength {
		return PadRight(v, maxLength, " ")
	}

	runes := []rune(v)
	if len(runes) > maxLength-3 {
		runes = runes[:maxLength-3]
	}
	return string(runes) + "..."
}

func PadRight(s string, length int, pad string) string {
	for len(s) < length {
		s = s + pad
	}
	return s
}

// AlignMultiline indents every line after the first so multiline values stay
// aligned with the start of the first line.
func AlignMultiline(value string, indent int) string {
	value = strings.ReplaceAll(value, "\r\n", "\n")
	value = strings.TrimRight(value, "\r\n")

	lines := strings.Split(value, "\n")
	pad := strings.Repeat(" ", indent)
	for i := 1; i < len(lines); i++ {
		lines[i] = pad + lines[i]
	}
	return strings.Join(lines, "\n")
}

const progressBarWidth = 10

func FormatProgressBar(percent int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := percent * progressBarWidth / 100
	return fmt.Sprintf("%3d%% %s", percent, strings.Repeat("■", filled))
}
