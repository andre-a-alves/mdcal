package calendar

import "time"

// Options represents the configuration for generating a calendar
type Options struct {
	Year             int
	Month            int
	FirstDayOfWeek   time.Weekday
	ShowCalendarWeek bool
	ShowWeekends     bool
	ShowComments     bool
	Justify          string
}

// NewOptions creates a new Options instance with default values
func NewOptions() Options {
	return Options{
		Year:             time.Now().Year(),
		Month:            0, // 0 means generate for the whole year
		FirstDayOfWeek:   time.Monday,
		ShowCalendarWeek: true,
		ShowWeekends:     true,
		ShowComments:     true,
		Justify:          "left",
	}
}
