package utils

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseWeekday(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected time.Weekday
	}{
		// Full names
		{
			name:     "Sunday full name",
			input:    "Sunday",
			expected: time.Sunday,
		},
		{
			name:     "Monday full name",
			input:    "Monday",
			expected: time.Monday,
		},
		{
			name:     "Tuesday full name",
			input:    "Tuesday",
			expected: time.Tuesday,
		},
		{
			name:     "Wednesday full name",
			input:    "Wednesday",
			expected: time.Wednesday,
		},
		{
			name:     "Thursday full name",
			input:    "Thursday",
			expected: time.Thursday,
		},
		{
			name:     "Friday full name",
			input:    "Friday",
			expected: time.Friday,
		},
		{
			name:     "Saturday full name",
			input:    "Saturday",
			expected: time.Saturday,
		},

		// Short names
		{
			name:     "Sunday short name",
			input:    "Sun",
			expected: time.Sunday,
		},
		{
			name:     "Monday short name",
			input:    "Mon",
			expected: time.Monday,
		},
		{
			name:     "Tuesday short name",
			input:    "Tue",
			expected: time.Tuesday,
		},
		{
			name:     "Tuesday alternative short name",
			input:    "Tues",
			expected: time.Tuesday,
		},
		{
			name:     "Wednesday short name",
			input:    "Wed",
			expected: time.Wednesday,
		},
		{
			name:     "Thursday short name",
			input:    "Thu",
			expected: time.Thursday,
		},
		{
			name:     "Thursday alternative short name 1",
			input:    "Thur",
			expected: time.Thursday,
		},
		{
			name:     "Thursday alternative short name 2",
			input:    "Thurs",
			expected: time.Thursday,
		},
		{
			name:     "Friday short name",
			input:    "Fri",
			expected: time.Friday,
		},
		{
			name:     "Saturday short name",
			input:    "Sat",
			expected: time.Saturday,
		},

		// Case insensitivity
		{
			name:     "Mixed case",
			input:    "MoNdAy",
			expected: time.Monday,
		},
		{
			name:     "All lowercase",
			input:    "tuesday",
			expected: time.Tuesday,
		},
		{
			name:     "All uppercase",
			input:    "WEDNESDAY",
			expected: time.Wednesday,
		},

		// Edge cases
		{
			name:     "Empty string",
			input:    "",
			expected: time.Monday, // Default
		},
		{
			name:     "Invalid day",
			input:    "NotADay",
			expected: time.Monday, // Default
		},
		{
			name:     "Partial match",
			input:    "Mond",
			expected: time.Monday, // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := ParseWeekday(tt.input)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("ParseWeekday(%q) mismatch (-want +got):\n%s", tt.input, diff)
			}
		})
	}
}
