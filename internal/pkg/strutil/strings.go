package strutil

import "unicode"

// CapitalizeFirst returns s with the first rune converted to uppercase.
func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
