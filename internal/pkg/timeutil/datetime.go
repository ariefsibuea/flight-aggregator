package timeutil

import (
	"fmt"
	"time"
)

const (
	ISO8601Basic      = "2006-01-02T15:04:05-0700"
	ISO8601NoTimezone = "2006-01-02T15:04:05"
)

// ParseDateTime converts a datetime string into a time.Time object.
func ParseDateTime(datetime string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		ISO8601Basic,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, datetime); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("failed to parse datetime %q", datetime)
}

// ParseDateTimeInLocation parses a datetime string with no timezone info using the given timezone name.
func ParseDateTimeInLocation(datetime, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, fmt.Errorf("unknown timezone %q: %w", timezone, err)
	}

	t, err := time.ParseInLocation(ISO8601NoTimezone, datetime, loc)
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot parse datetime %q: %w", datetime, err)
	}

	return t, nil
}
