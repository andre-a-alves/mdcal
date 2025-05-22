package calendar

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestGenerateCalendarHeader(t *testing.T) {
	tests := []struct {
		name     string
		year     int
		month    time.Month
		expected string
	}{
		{
			name:     "January 2023",
			year:     2023,
			month:    time.January,
			expected: "# January 2023\n\n",
		},
		{
			name:     "December 2025",
			year:     2025,
			month:    time.December,
			expected: "# December 2025\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := generateCalendarHeader(tt.year, tt.month)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("generateCalendarHeader(%d, %v) mismatch (-want +got):\n%s", tt.year, tt.month, diff)
			}
		})
	}
}

func TestGetWeekdayNames(t *testing.T) {
	tests := []struct {
		name           string
		firstDayOfWeek time.Weekday
		showWeekends   bool
		expectedShort  []string
		expectedFull   []string
	}{
		{
			name:           "Monday first, with weekends",
			firstDayOfWeek: time.Monday,
			showWeekends:   true,
			expectedShort:  []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			expectedFull:   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"},
		},
		{
			name:           "Sunday first, with weekends",
			firstDayOfWeek: time.Sunday,
			showWeekends:   true,
			expectedShort:  []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"},
			expectedFull:   []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
		},
		{
			name:           "Wednesday first, with weekends",
			firstDayOfWeek: time.Wednesday,
			showWeekends:   true,
			expectedShort:  []string{"Wed", "Thu", "Fri", "Sat", "Sun", "Mon", "Tue"},
			expectedFull:   []string{"Wednesday", "Thursday", "Friday", "Saturday", "Sunday", "Monday", "Tuesday"},
		},
		{
			name:           "Monday first, without weekends",
			firstDayOfWeek: time.Monday,
			showWeekends:   false,
			expectedShort:  []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			expectedFull:   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
		{
			name:           "Sunday first, without weekends",
			firstDayOfWeek: time.Sunday,
			showWeekends:   false,
			expectedShort:  []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			expectedFull:   []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shortNames, fullNames := getWeekdayNames(tt.firstDayOfWeek, tt.showWeekends)

			if diff := cmp.Diff(tt.expectedShort, shortNames); diff != "" {
				t.Errorf("getWeekdayNames() short names mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.expectedFull, fullNames); diff != "" {
				t.Errorf("getWeekdayNames() full names mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPrepareColumnHeaders(t *testing.T) {
	tests := []struct {
		name             string
		dayShortNames    []string
		dayFullNames     []string
		useShortDayNames bool
		showCalendarWeek bool
		showComments     bool
		justify          string
		expectedHeaders  []string
		expectedWidths   []int
	}{
		{
			name:             "Short names, with week numbers and comments",
			dayShortNames:    []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			dayFullNames:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
			useShortDayNames: true,
			showCalendarWeek: true,
			showComments:     true,
			justify:          "left",
			expectedHeaders:  []string{"CW  ", "Mon", "Tue", "Wed", "Thu", "Fri", "Comments"},
			expectedWidths:   []int{4, 3, 3, 3, 3, 3, 8},
		},
		{
			name:             "Full names, with week numbers and comments",
			dayShortNames:    []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			dayFullNames:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
			useShortDayNames: false,
			showCalendarWeek: true,
			showComments:     true,
			justify:          "left",
			expectedHeaders:  []string{"CW  ", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Comments"},
			expectedWidths:   []int{4, 6, 7, 9, 8, 6, 8},
		},
		{
			name:             "Short names, no week numbers, with comments",
			dayShortNames:    []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			dayFullNames:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
			useShortDayNames: true,
			showCalendarWeek: false,
			showComments:     true,
			justify:          "left",
			expectedHeaders:  []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Comments"},
			expectedWidths:   []int{3, 3, 3, 3, 3, 8},
		},
		{
			name:             "Short names, with week numbers, no comments",
			dayShortNames:    []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			dayFullNames:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
			useShortDayNames: true,
			showCalendarWeek: true,
			showComments:     false,
			justify:          "left",
			expectedHeaders:  []string{"CW  ", "Mon", "Tue", "Wed", "Thu", "Fri"},
			expectedWidths:   []int{4, 3, 3, 3, 3, 3},
		},
		{
			name:             "Center justify",
			dayShortNames:    []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			dayFullNames:     []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"},
			useShortDayNames: true,
			showCalendarWeek: true,
			showComments:     true,
			justify:          "center",
			expectedHeaders:  []string{"CW   ", "Mon", "Tue", "Wed", "Thu", "Fri", "Comments"},
			expectedWidths:   []int{5, 3, 3, 3, 3, 3, 8},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers, widths := prepareColumnHeaders(tt.dayShortNames, tt.dayFullNames, tt.useShortDayNames,
				tt.showCalendarWeek, tt.showComments, tt.justify)

			if diff := cmp.Diff(tt.expectedHeaders, headers); diff != "" {
				t.Errorf("prepareColumnHeaders() headers mismatch (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tt.expectedWidths, widths); diff != "" {
				t.Errorf("prepareColumnHeaders() widths mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGenerateTableHeader(t *testing.T) {
	tests := []struct {
		name          string
		columnHeaders []string
		columnWidths  []int
		justify       string
		expected      string
	}{
		{
			name:          "Left justify",
			columnHeaders: []string{"CW", "Mon", "Tue", "Wed"},
			columnWidths:  []int{2, 3, 3, 3},
			justify:       "left",
			expected:      "| CW | Mon | Tue | Wed |\n| :- | :-- | :-- | :-- |\n",
		},
		{
			name:          "Center justify",
			columnHeaders: []string{"CW", "Mon", "Tue", "Wed"},
			columnWidths:  []int{2, 3, 3, 3},
			justify:       "center",
			expected:      "| CW | Mon | Tue | Wed |\n| :-: | :-: | :-: | :-: |\n",
		},
		{
			name:          "Right justify",
			columnHeaders: []string{"CW", "Mon", "Tue", "Wed"},
			columnWidths:  []int{2, 3, 3, 3},
			justify:       "right",
			expected:      "| CW | Mon | Tue | Wed |\n| -: | --: | --: | --: |\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := generateTableHeader(tt.columnHeaders, tt.columnWidths, tt.justify)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("generateTableHeader() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestConvertToWeekdays(t *testing.T) {
	tests := []struct {
		name          string
		dayShortNames []string
		expected      []time.Weekday
	}{
		{
			name:          "All days",
			dayShortNames: []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"},
			expected:      []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday},
		},
		{
			name:          "Weekdays only",
			dayShortNames: []string{"Mon", "Tue", "Wed", "Thu", "Fri"},
			expected:      []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
		},
		{
			name:          "Custom order",
			dayShortNames: []string{"Wed", "Thu", "Fri", "Mon", "Tue"},
			expected:      []time.Weekday{time.Wednesday, time.Thursday, time.Friday, time.Monday, time.Tuesday},
		},
		{
			name:          "Empty list",
			dayShortNames: []string{},
			expected:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := convertToWeekdays(tt.dayShortNames)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("convertToWeekdays() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestCalculateMonthBoundaries(t *testing.T) {
	tests := []struct {
		name                 string
		year                 int
		month                time.Month
		firstDayOfWeek       time.Weekday
		expectedFirstOfMonth time.Time
		expectedLastOfMonth  time.Time
		expectedWeekStart    time.Time
	}{
		{
			name:                 "January 2023, Monday first",
			year:                 2023,
			month:                time.January,
			firstDayOfWeek:       time.Monday,
			expectedFirstOfMonth: time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC),
			expectedLastOfMonth:  time.Date(2023, time.January, 31, 0, 0, 0, 0, time.UTC),
			expectedWeekStart:    time.Date(2022, time.December, 26, 0, 0, 0, 0, time.UTC),
		},
		{
			name:                 "February 2024 (leap year), Sunday first",
			year:                 2024,
			month:                time.February,
			firstDayOfWeek:       time.Sunday,
			expectedFirstOfMonth: time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
			expectedLastOfMonth:  time.Date(2024, time.February, 29, 0, 0, 0, 0, time.UTC),
			expectedWeekStart:    time.Date(2024, time.January, 28, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			firstOfMonth, lastOfMonth, weekStart := calculateMonthBoundaries(tt.year, tt.month, tt.firstDayOfWeek)

			if !firstOfMonth.Equal(tt.expectedFirstOfMonth) {
				t.Errorf("calculateMonthBoundaries() firstOfMonth = %v, want %v",
					firstOfMonth, tt.expectedFirstOfMonth)
			}

			if !lastOfMonth.Equal(tt.expectedLastOfMonth) {
				t.Errorf("calculateMonthBoundaries() lastOfMonth = %v, want %v",
					lastOfMonth, tt.expectedLastOfMonth)
			}

			if !weekStart.Equal(tt.expectedWeekStart) {
				t.Errorf("calculateMonthBoundaries() weekStart = %v, want %v",
					weekStart, tt.expectedWeekStart)
			}
		})
	}
}

func TestGenerateWeekRow(t *testing.T) {
	tests := []struct {
		name             string
		cur              time.Time
		month            time.Month
		weekDays         []time.Weekday
		firstDayOfWeek   time.Weekday
		columnWidths     []int
		showCalendarWeek bool
		showComments     bool
		expected         string
	}{
		{
			name:             "First week of January 2023, Monday first",
			cur:              time.Date(2022, time.December, 26, 0, 0, 0, 0, time.UTC),
			month:            time.January,
			weekDays:         []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			firstDayOfWeek:   time.Monday,
			columnWidths:     []int{4, 3, 3, 3, 3, 3, 8}, // Need 7 elements: 1 for week number, 5 for weekdays, 1 for comments
			showCalendarWeek: true,
			showComments:     true,
			expected:         "| _52_ |     |     |     |     |     |          |\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := generateWeekRow(tt.cur, tt.month, tt.weekDays, tt.firstDayOfWeek,
				tt.columnWidths, tt.showCalendarWeek, tt.showComments)
			if diff := cmp.Diff(tt.expected, actual); diff != "" {
				t.Errorf("generateWeekRow() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestValidateDateRange(t *testing.T) {
	tests := []struct {
		name          string
		options       Options
		expectedValid bool
		expectedMsg   string
	}{
		{
			name: "No end date specified",
			options: Options{
				Year:  2023,
				Month: intPtr(5),
			},
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name: "Valid date range (same year)",
			options: Options{
				Year:     2023,
				Month:    intPtr(5),
				EndYear:  intPtr(2023),
				EndMonth: intPtr(7),
			},
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name: "Valid date range (different years)",
			options: Options{
				Year:     2023,
				Month:    intPtr(11),
				EndYear:  intPtr(2024),
				EndMonth: intPtr(2),
			},
			expectedValid: true,
			expectedMsg:   "",
		},
		{
			name: "Invalid date range (end before start)",
			options: Options{
				Year:     2023,
				Month:    intPtr(5),
				EndYear:  intPtr(2023),
				EndMonth: intPtr(3),
			},
			expectedValid: false,
			expectedMsg:   "Error: End date cannot be before start date\n",
		},
		{
			name: "Invalid date range (different years)",
			options: Options{
				Year:     2023,
				Month:    intPtr(5),
				EndYear:  intPtr(2022),
				EndMonth: intPtr(7),
			},
			expectedValid: false,
			expectedMsg:   "Error: End date cannot be before start date\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := validateDateRange(tt.options)

			if valid != tt.expectedValid {
				t.Errorf("validateDateRange() valid = %v, want %v", valid, tt.expectedValid)
			}

			if msg != tt.expectedMsg {
				t.Errorf("validateDateRange() msg = %q, want %q", msg, tt.expectedMsg)
			}
		})
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
