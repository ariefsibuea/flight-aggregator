package timezone

import "time"

// Indonesian timezone offsets
const (
	WIB  = "+07:00" //  Western Indonesian Time is UTC+7
	WITA = "+08:00" //  Central Indonesian Time is UTC+8
	WIT  = "+09:00" //  Eastern Indonesian Time is UTC+9
)

// LocationByName returns the time.Location for a given timezone name
func LocationByName(name string) (*time.Location, error) {
	return time.LoadLocation(name)
}
