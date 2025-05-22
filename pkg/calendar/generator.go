package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/andre-a-alves/mdcal/pkg/utils"
)

// GenerateMonthCalendar generates a markdown calendar for the specified month
func GenerateMonthCalendar(options Options) string {
	var sb strings.Builder

	// header
	month := time.Month(options.Month)
	sb.WriteString(fmt.Sprintf("# %s %d\n\n", month.String(), options.Year))

	// all weekday names with full names
	allShortNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	allFullNames := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

	// convert time.Weekday (Mon=1…) to 0-based Monday=0…Sunday=6
	firstIndex := (int(options.FirstDayOfWeek) + 6) % 7
	shiftedShort := append(allShortNames[firstIndex:], allShortNames[:firstIndex]...)
	shiftedFull := append(allFullNames[firstIndex:], allFullNames[:firstIndex]...)

	// strip weekends if needed
	var dayShortNames, dayFullNames []string
	for i, d := range shiftedShort {
		if !options.ShowWeekends && (d == "Sat" || d == "Sun") {
			continue
		}
		dayShortNames = append(dayShortNames, d)
		dayFullNames = append(dayFullNames, shiftedFull[i])
	}

	// prepare column headers and widths
	var columnHeaders []string
	var columnWidths []int
	if options.ShowCalendarWeek {
		if strings.ToLower(options.Justify) == "center" {
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
	if options.ShowComments {
		columnHeaders = append(columnHeaders, "Comments")
		columnWidths = append(columnWidths, len("Comments"))
	}

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
		sep := utils.SeparatorCell(w, options.Justify)
		sb.WriteString(" " + sep + " |")
	}
	sb.WriteString("\n")

	// map back to weekdays
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

	// compute first and last days
	firstOfMonth := time.Date(options.Year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := time.Date(options.Year, month+1, 0, 0, 0, 0, 0, time.UTC)
	shift := (int(firstOfMonth.Weekday()) - int(options.FirstDayOfWeek) + 7) % 7
	weekStart := firstOfMonth.AddDate(0, 0, -shift)

	// print each week
	for cur := weekStart; !cur.After(lastOfMonth); cur = cur.AddDate(0, 0, 7) {
		var cells []string
		if options.ShowCalendarWeek {
			_, w := cur.ISOWeek()
			// leave calendar week right‐justified if you prefer, or treat like a header
			cells = append(cells, fmt.Sprintf("_%d_", w))
		}
		for _, wd := range weekDays {
			delta := (int(wd) - int(options.FirstDayOfWeek) + 7) % 7
			cd := cur.AddDate(0, 0, delta)
			if cd.Month() == month {
				// no leading zero, left-justify
				cells = append(cells, fmt.Sprintf("%d", cd.Day()))
			} else {
				cells = append(cells, "")
			}
		}
		if options.ShowComments {
			cells = append(cells, "")
		}

		sb.WriteString("|")
		for i, cell := range cells {
			w := columnWidths[i]
			// left-justify all cells now
			sb.WriteString(" " + utils.PadRight(cell, w) + " |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// PrintCalendar generates and returns the calendar based on the provided options
func PrintCalendar(options Options) string {
	var result strings.Builder

	// Validate date range the end date is specified
	if options.EndYear != 0 && options.EndMonth != 0 {
		startDate := time.Date(options.Year, time.Month(options.Month), 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(options.EndYear, time.Month(options.EndMonth), 1, 0, 0, 0, 0, time.UTC)

		if startDate.After(endDate) {
			return "Error: End date cannot be before start date\n"
		}
	}

	if options.Month == 0 {
		// Generate calendar for the whole year
		for m := 1; m <= 12; m++ {
			optionsCopy := options
			optionsCopy.Month = m
			result.WriteString(GenerateMonthCalendar(optionsCopy))
			result.WriteString("\n") // Add a blank line between months
		}
	} else if options.EndYear != 0 && options.EndMonth != 0 {
		// Generate calendar for a range of months
		currentYear := options.Year
		currentMonth := options.Month

		for {
			optionsCopy := options
			optionsCopy.Year = currentYear
			optionsCopy.Month = currentMonth
			result.WriteString(GenerateMonthCalendar(optionsCopy))
			result.WriteString("\n") // Add a blank line between months

			// Move to the next month
			currentMonth++
			if currentMonth > 12 {
				currentMonth = 1
				currentYear++
			}

			// Check if we've reached the end of the range
			if currentYear > options.EndYear || (currentYear == options.EndYear && currentMonth > options.EndMonth) {
				break
			}
		}
	} else {
		// Generate calendar for the specific month
		result.WriteString(GenerateMonthCalendar(options))
	}

	return result.String()
}
