package calendar

import "time"

// Options represents the configuration for generating a calendar
type Options struct {
	Year             int
	Month            int
	EndYear          int // End year for range
	EndMonth         int // End month for range
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
		EndYear:          0, // 0 means no end year specified
		EndMonth:         0, // 0 means no end month specified
		FirstDayOfWeek:   time.Monday,
		ShowCalendarWeek: true,
		ShowWeekends:     true,
		ShowComments:     true,
		Justify:          "left",
	}
}
