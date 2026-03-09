package strutil

import (
	"strconv"
	"strings"
	"unicode"
)

// CapitalizeFirst returns s with the first rune converted to uppercase.
func CapitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func FormatCurrency(amount int64, symbol string) string {
	isNegative := amount < 0
	if isNegative {
		amount = -amount
	}

	s := strconv.FormatInt(amount, 10)
	n := len(s)

	var b strings.Builder
	b.WriteString(symbol)
	b.WriteString(" ")

	if isNegative {
		b.WriteByte('-')
	}

	for i, c := range s {
		if i > 0 && (n-i)%3 == 0 {
			b.WriteByte('.')
		}
		b.WriteRune(c)
	}

	return b.String()
}
