package main

import "unicode/utf8"

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
