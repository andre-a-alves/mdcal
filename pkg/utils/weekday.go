package utils

import (
	"strings"
	"time"
)

// ParseWeekday converts a string representation of a weekday to time.Weekday
func ParseWeekday(day string) time.Weekday {
	switch strings.ToLower(day) {
	case "sunday":
		return time.Sunday
	case "monday":
		return time.Monday
	case "tuesday":
		return time.Tuesday
	case "wednesday":
		return time.Wednesday
	case "thursday":
		return time.Thursday
	case "friday":
		return time.Friday
	case "saturday":
		return time.Saturday
	default:
		return time.Monday // Default to Monday
	}
}
