package timezone

import "time"

// Indonesian timezone offsets
const (
	WIB  = "+07:00" //  Wester is UTC+7
	WITA = "+08:00" //  Central is UTC+8
	WIT  = "+09:00" //  Eastern is UTC+9
)

// Location returns the time.Location for a given timezone name
func Location(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}

// OffsetDuration returns the time.Duration for a timezone offset string
func OffsetDuration(offset string) (time.Duration, error) {
	return time.ParseDuration(offset)
}
