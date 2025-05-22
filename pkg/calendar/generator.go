package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/andre-a-alves/mdcal/pkg/utils"
)

// generateCalendarHeader creates the header for the calendar with month and year
func generateCalendarHeader(year int, month time.Month) string {
	return fmt.Sprintf("# %s %d\n\n", month.String(), year)
}

// getWeekdayNames returns the short and full names of weekdays based on the first day of the week
func getWeekdayNames(firstDayOfWeek time.Weekday, showWeekends bool) ([]string, []string) {
	allShortNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	allFullNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	// convert time.Weekday (Mon=1…) to 0-based Monday=0…Sunday=6
	firstIndex := (int(firstDayOfWeek) + 6) % 7
	shiftedShort := append(allShortNames[firstIndex:], allShortNames[:firstIndex]...)
	shiftedFull := append(allFullNames[firstIndex:], allFullNames[:firstIndex]...)

	// strip weekends if needed
	var dayShortNames, dayFullNames []string
	for i, d := range shiftedShort {
		if !showWeekends && (d == "Sat" || d == "Sun") {
			continue
		}
		dayShortNames = append(dayShortNames, d)
		dayFullNames = append(dayFullNames, shiftedFull[i])
	}

	return dayShortNames, dayFullNames
}

// prepareColumnHeaders creates the column headers and their widths
func prepareColumnHeaders(dayFullNames []string, showCalendarWeek bool, showComments bool, justify string) ([]string, []int) {
	var columnHeaders []string
	var columnWidths []int

	if showCalendarWeek {
		if strings.ToLower(justify) == "center" {
			columnHeaders = append(columnHeaders, "CW   ")
			columnWidths = append(columnWidths, len("CW   "))
		} else {
			columnHeaders = append(columnHeaders, "CW  ")
			columnWidths = append(columnWidths, len("CW  "))
		}
	}

	for _, d := range dayFullNames {
		columnHeaders = append(columnHeaders, d)
		columnWidths = append(columnWidths, len(d))
	}

	if showComments {
		columnHeaders = append(columnHeaders, "Comments")
		columnWidths = append(columnWidths, len("Comments"))
	}

	return columnHeaders, columnWidths
}

// generateTableHeader creates the header row and separator row for the markdown table
func generateTableHeader(columnHeaders []string, columnWidths []int, justify string) string {
	var sb strings.Builder

	// header row
	sb.WriteString("|")
	for i, h := range columnHeaders {
		w := columnWidths[i]
		sb.WriteString(" " + utils.PadRight(h, w) + " |")
	}
	sb.WriteString("\n")

	// separator row
	sb.WriteString("|")
	for _, w := range columnWidths {
		sep := utils.SeparatorCell(w, justify)
		sb.WriteString(" " + sep + " |")
	}
	sb.WriteString("\n")

	return sb.String()
}

// convertToWeekdays converts short day names to time.Weekday values
func convertToWeekdays(dayShortNames []string) []time.Weekday {
	var weekDays []time.Weekday
	for _, d := range dayShortNames {
		switch d {
		case "Mon":
			weekDays = append(weekDays, time.Monday)
		case "Tue":
			weekDays = append(weekDays, time.Tuesday)
		case "Wed":
			weekDays = append(weekDays, time.Wednesday)
		case "Thu":
			weekDays = append(weekDays, time.Thursday)
		case "Fri":
			weekDays = append(weekDays, time.Friday)
		case "Sat":
			weekDays = append(weekDays, time.Saturday)
		case "Sun":
			weekDays = append(weekDays, time.Sunday)
		}
	}
	return weekDays
}

// calculateMonthBoundaries computes the first and last days of the month and the start of the first week
func calculateMonthBoundaries(year int, month time.Month, firstDayOfWeek time.Weekday) (firstOfMonth time.Time, lastOfMonth time.Time, weekStart time.Time) {
	firstOfMonth = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth = time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	shift := (int(firstOfMonth.Weekday()) - int(firstDayOfWeek) + 7) % 7
	weekStart = firstOfMonth.AddDate(0, 0, -shift)

	return firstOfMonth, lastOfMonth, weekStart
}

// generateWeekRow creates a single week row for the calendar
func generateWeekRow(cur time.Time, month time.Month, weekDays []time.Weekday, firstDayOfWeek time.Weekday,
	columnWidths []int, showCalendarWeek bool, showComments bool) string {
	var sb strings.Builder
	var cells []string

	if showCalendarWeek {
		_, w := cur.ISOWeek()
		cells = append(cells, fmt.Sprintf("_%d_", w))
	}

	for _, wd := range weekDays {
		delta := (int(wd) - int(firstDayOfWeek) + 7) % 7
		cd := cur.AddDate(0, 0, delta)
		if cd.Month() == month {
			cells = append(cells, fmt.Sprintf("%d", cd.Day()))
		} else {
			cells = append(cells, "")
		}
	}

	if showComments {
		cells = append(cells, "")
	}

	sb.WriteString("|")
	for i, cell := range cells {
		w := columnWidths[i]
		sb.WriteString(" " + utils.PadRight(cell, w) + " |")
	}
	sb.WriteString("\n")

	return sb.String()
}

// GenerateMonthCalendar generates a markdown calendar for the specified month
func GenerateMonthCalendar(options Options) string {
	var sb strings.Builder
	month := time.Month(*options.Month)

	// Add calendar header
	sb.WriteString(generateCalendarHeader(options.Year, month))

	// Get weekday names
	dayShortNames, dayFullNames := getWeekdayNames(options.FirstDayOfWeek, options.ShowWeekends)

	// Prepare column headers and widths
	columnHeaders, columnWidths := prepareColumnHeaders(dayFullNames, options.ShowCalendarWeek,
		options.ShowComments, options.Justify)

	// Generate table header
	sb.WriteString(generateTableHeader(columnHeaders, columnWidths, options.Justify))

	// Convert short day names to weekdays
	weekDays := convertToWeekdays(dayShortNames)

	// Calculate month boundaries
	_, lastOfMonth, weekStart := calculateMonthBoundaries(options.Year, month, options.FirstDayOfWeek)

	// Generate each week row
	for cur := weekStart; !cur.After(lastOfMonth); cur = cur.AddDate(0, 0, 7) {
		sb.WriteString(generateWeekRow(cur, month, weekDays, options.FirstDayOfWeek,
			columnWidths, options.ShowCalendarWeek, options.ShowComments))
	}

	return sb.String()
}

// validateDateRange checks if the end date is after the start date
func validateDateRange(options Options) (bool, string) {
	if options.EndYear == nil || options.EndMonth == nil {
		return true, ""
	}

	monthValue := 1
	if options.Month != nil {
		monthValue = *options.Month
	}

	startDate := time.Date(options.Year, time.Month(monthValue), 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(*options.EndYear, time.Month(*options.EndMonth), 1, 0, 0, 0, 0, time.UTC)

	if startDate.After(endDate) {
		return false, "Error: End date cannot be before start date\n"
	}

	return true, ""
}

// generateYearCalendar creates a calendar for the entire year
func generateYearCalendar(options Options) string {
	var result strings.Builder

	for m := 1; m <= 12; m++ {
		optionsCopy := options
		monthValue := m
		optionsCopy.Month = &monthValue
		result.WriteString(GenerateMonthCalendar(optionsCopy))
		result.WriteString("\n") // Add a blank line between months
	}

	return result.String()
}

// generateMonthRangeCalendar creates a calendar for a range of months
func generateMonthRangeCalendar(options Options) string {
	var result strings.Builder
	currentYear := options.Year
	currentMonth := *options.Month

	for {
		optionsCopy := options
		optionsCopy.Year = currentYear
		monthValue := currentMonth
		optionsCopy.Month = &monthValue
		result.WriteString(GenerateMonthCalendar(optionsCopy))
		result.WriteString("\n") // Add a blank line between months

		// Move to the next month
		currentMonth++
		if currentMonth > 12 {
			currentMonth = 1
			currentYear++
		}

		// Check if we've reached the end of the range
		if currentYear > *options.EndYear || (currentYear == *options.EndYear && currentMonth > *options.EndMonth) {
			break
		}
	}

	return result.String()
}

// PrintCalendar generates and returns the calendar based on the provided options
func PrintCalendar(options Options) string {
	// Validate date range if the end date is specified
	valid, errorMsg := validateDateRange(options)
	if !valid {
		return errorMsg
	}

	if options.Month == nil {
		// Generate calendar for the whole year
		return generateYearCalendar(options)
	} else if options.EndYear != nil && options.EndMonth != nil {
		// Generate calendar for a range of months
		return generateMonthRangeCalendar(options)
	} else {
		// Generate calendar for the specific month
		return GenerateMonthCalendar(options)
	}
}
