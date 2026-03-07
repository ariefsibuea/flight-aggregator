package timeutil

import (
	"fmt"
	"strings"
)

// ParseDuration parses a duration string in "Xh Ym" format and returns the total minutes.
// Supports "Xh Ym", "Xh", and "Xm" variants.
func ParseDuration(duration string) (int, error) {
	duration = strings.TrimSpace(duration)
	var h, m int

	if _, err := fmt.Sscanf(duration, "%dh %dm", &h, &m); err == nil {
		return h*60 + m, nil
	}
	if _, err := fmt.Sscanf(duration, "%dh", &h); err == nil {
		return h * 60, nil
	}
	if _, err := fmt.Sscanf(duration, "%dm", &m); err == nil {
		return m, nil
	}

	return 0, fmt.Errorf("cannot parse duration %q", duration)
}

// FormatDuration formats a duration in minutes as "Xh Ym".
func FormatDuration(minutes int) string {
	h := minutes / 60
	m := minutes % 60
	return fmt.Sprintf("%dh %dm", h, m)
}
