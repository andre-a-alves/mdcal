package interactive

import (
	"bufio"
	"strings"
	"testing"
	"time"

	"github.com/andre-a-alves/mdcal/pkg/calendar"
	"github.com/google/go-cmp/cmp"
)

func TestReadInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Empty string",
			input:    "\n",
			expected: "",
		},
		{
			name:     "String with whitespace",
			input:    "  hello  \n",
			expected: "hello",
		},
		{
			name:     "String without whitespace",
			input:    "hello\n",
			expected: "hello",
		},
		{
			name:     "Multiple lines (only reads first)",
			input:    "line1\nline2\n",
			expected: "line1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			actual := readInput(reader)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("readInput() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPromptForYear(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		initialYear  int
		expectedYear int
	}{
		{
			name:         "Empty input (use default)",
			input:        "\n",
			initialYear:  2023,
			expectedYear: 2023,
		},
		{
			name:         "Valid year",
			input:        "2025\n",
			initialYear:  2023,
			expectedYear: 2025,
		},
		{
			name:         "Invalid year (use default)",
			input:        "not-a-year\n",
			initialYear:  2023,
			expectedYear: 2023,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			options := calendar.Options{Year: tt.initialYear}

			promptForYear(reader, &options)

			if diff := cmp.Diff(tt.expectedYear, options.Year); diff != "" {
				t.Errorf("promptForYear() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPromptForMonth(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		initialMonth  *int
		expectedMonth *int
	}{
		{
			name:          "Empty input (whole year)",
			input:         "\n",
			initialMonth:  intPtr(5),
			expectedMonth: nil,
		},
		{
			name:          "Valid month",
			input:         "7\n",
			initialMonth:  intPtr(5),
			expectedMonth: intPtr(7),
		},
		{
			name:          "Invalid month (use whole year)",
			input:         "not-a-month\n",
			initialMonth:  intPtr(5),
			expectedMonth: nil,
		},
		{
			name:          "Month out of range (use whole year)",
			input:         "13\n",
			initialMonth:  intPtr(5),
			expectedMonth: nil,
		},
		{
			name:          "Month zero (use whole year)",
			input:         "0\n",
			initialMonth:  intPtr(5),
			expectedMonth: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			options := calendar.Options{Month: tt.initialMonth}

			promptForMonth(reader, &options)

			// Compare month values, handling nil pointers
			if (tt.expectedMonth == nil && options.Month != nil) ||
				(tt.expectedMonth != nil && options.Month == nil) {
				t.Errorf("promptForMonth() expected month: %v, got: %v",
					ptrValue(tt.expectedMonth), ptrValue(options.Month))
			} else if tt.expectedMonth != nil && options.Month != nil &&
				*tt.expectedMonth != *options.Month {
				t.Errorf("promptForMonth() expected month: %d, got: %d",
					*tt.expectedMonth, *options.Month)
			}
		})
	}
}

func TestPromptForWeekStartDay(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		initialFirstDay  time.Weekday
		expectedFirstDay time.Weekday
	}{
		{
			name:             "Empty input (use default)",
			input:            "\n",
			initialFirstDay:  time.Monday,
			expectedFirstDay: time.Monday,
		},
		{
			name:             "Valid day (full name)",
			input:            "Sunday\n",
			initialFirstDay:  time.Monday,
			expectedFirstDay: time.Sunday,
		},
		{
			name:             "Valid day (short name)",
			input:            "Fri\n",
			initialFirstDay:  time.Monday,
			expectedFirstDay: time.Friday,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			options := calendar.Options{FirstDayOfWeek: tt.initialFirstDay}

			promptForWeekStartDay(reader, &options)

			if diff := cmp.Diff(tt.expectedFirstDay, options.FirstDayOfWeek); diff != "" {
				t.Errorf("promptForWeekStartDay() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPromptForBooleanOption(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Empty input (default: false)",
			input:    "\n",
			expected: true,
		},
		{
			name:     "Yes (lowercase)",
			input:    "yes\n",
			expected: false,
		},
		{
			name:     "Y (lowercase)",
			input:    "y\n",
			expected: false,
		},
		{
			name:     "Yes (uppercase)",
			input:    "YES\n",
			expected: false,
		},
		{
			name:     "No (lowercase)",
			input:    "no\n",
			expected: true,
		},
		{
			name:     "N (lowercase)",
			input:    "n\n",
			expected: true,
		},
		{
			name:     "Invalid input (default: false)",
			input:    "invalid\n",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			actual := promptForBooleanOption(reader, "Prompt: ")

			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("promptForBooleanOption() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPromptForJustification(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		initialJustify  string
		expectedJustify string
	}{
		{
			name:            "Empty input (use default)",
			input:           "\n",
			initialJustify:  "left",
			expectedJustify: "left",
		},
		{
			name:            "Valid justify (left)",
			input:           "left\n",
			initialJustify:  "center",
			expectedJustify: "left",
		},
		{
			name:            "Valid justify (center)",
			input:           "center\n",
			initialJustify:  "left",
			expectedJustify: "center",
		},
		{
			name:            "Valid justify (right)",
			input:           "right\n",
			initialJustify:  "left",
			expectedJustify: "right",
		},
		{
			name:            "Invalid justify (use default)",
			input:           "invalid\n",
			initialJustify:  "left",
			expectedJustify: "left",
		},
		{
			name:            "Mixed case justify",
			input:           "CeNtEr\n",
			initialJustify:  "left",
			expectedJustify: "center",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := bufio.NewReader(strings.NewReader(tt.input))
			options := calendar.Options{Justify: tt.initialJustify}

			promptForJustification(reader, &options)

			if diff := cmp.Diff(tt.expectedJustify, options.Justify); diff != "" {
				t.Errorf("promptForJustification() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func ptrValue(p *int) interface{} {
	if p == nil {
		return nil
	}
	return *p
}
