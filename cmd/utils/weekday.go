package utils

import (
	"strings"
	"time"
)

// ParseWeekday converts a string representation of a weekday to time.Weekday
func ParseWeekday(day string) time.Weekday {
	day = strings.ToLower(day)
	switch {
	case day == "sunday" || day == "sun":
		return time.Sunday
	case day == "monday" || day == "mon":
		return time.Monday
	case day == "tuesday" || day == "tue" || day == "tues":
		return time.Tuesday
	case day == "wednesday" || day == "wed":
		return time.Wednesday
	case day == "thursday" || day == "thu" || day == "thur" || day == "thurs":
		return time.Thursday
	case day == "friday" || day == "fri":
		return time.Friday
	case day == "saturday" || day == "sat":
		return time.Saturday
	default:
		return time.Monday // Default to Monday
	}
}
