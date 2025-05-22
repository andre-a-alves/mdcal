package calendar

import "time"

// Options represents the configuration for generating a calendar
type Options struct {
	Year             int
	Month            *int
	EndYear          *int // End year for range
	EndMonth         *int // End month for range
	FirstDayOfWeek   time.Weekday
	ShowCalendarWeek bool
	ShowWeekends     bool
	ShowComments     bool
	UseShortDayNames bool // Use short day names (Mon, Tue, etc.) instead of full names
	Justify          string
}

// NewOptions creates a new Options instance with default values
func NewOptions() Options {
	return Options{
		Year:             time.Now().Year(),
		Month:            nil, // nil means generate for the whole year
		EndYear:          nil, // nil means no end year specified
		EndMonth:         nil, // nil means no end month specified
		FirstDayOfWeek:   time.Monday,
		ShowCalendarWeek: true,
		ShowWeekends:     true,
		ShowComments:     true,
		UseShortDayNames: false,
		Justify:          "left",
	}
}
