package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type CalendarOptions struct {
	Year             int
	Month            int
	FirstDayOfWeek   time.Weekday
	ShowCalendarWeek bool
	ShowWeekends     bool
	ShowComments     bool
	Justify          string
}

func main() {
	// Define flags
	weekStart := flag.String("weekstart", "monday", "First day of the week (monday, sunday, etc.)")
	showCalWeek := flag.Bool("week", true, "Show calendar week numbers")
	workweek := flag.Bool("workweek", false, "Show only workdays (Monday-Friday)")
	showComments := flag.Bool("comments", true, "Add a comments column")
	versionFlag := flag.Bool("version", false, "Print version information")
	justify := flag.String("justify", "left", "Cell justification: left, center, or right")

	// Parse flags
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println("mdcal v0.1.0")
		os.Exit(0)
	}

	// Initialize options with defaults
	options := CalendarOptions{
		Year:             time.Now().Year(),
		Month:            0, // 0 means generate for the whole year
		FirstDayOfWeek:   parseWeekday(*weekStart),
		ShowCalendarWeek: *showCalWeek,
		ShowWeekends:     !*workweek,
		ShowComments:     *showComments,
		Justify:          *justify,
	}

	// Check if any flags or arguments were provided
	noFlagsOrArgs := flag.NFlag() == 0 && len(flag.Args()) == 0

	if noFlagsOrArgs {
		// Run in interactive mode
		runInteractiveMode(&options)
	} else {
		// Handle unnamed arguments (year and month)
		args := flag.Args()
		if len(args) > 0 {
			if year, err := strconv.Atoi(args[0]); err == nil {
				options.Year = year
			} else {
				fmt.Println("Invalid year, using current year")
			}
		}

		if len(args) > 1 {
			if month, err := strconv.Atoi(args[1]); err == nil && month >= 1 && month <= 12 {
				options.Month = month
			} else {
				fmt.Println("Invalid month, generating calendar for the whole year")
			}
		}
	}

	// Generate and print calendar
	printCalendar(options)
}

func runInteractiveMode(options *CalendarOptions) {
	reader := bufio.NewReader(os.Stdin)

	// Get year
	fmt.Printf("Enter year (default: %d): ", options.Year)
	yearInput, _ := reader.ReadString('\n')
	yearInput = strings.TrimSpace(yearInput)
	if yearInput != "" {
		if year, err := strconv.Atoi(yearInput); err == nil {
			options.Year = year
		} else {
			fmt.Println("Invalid year, using current year")
		}
	}

	// Get month
	fmt.Print("Enter month (1-12, empty for whole year): ")
	monthInput, _ := reader.ReadString('\n')
	monthInput = strings.TrimSpace(monthInput)
	if monthInput != "" {
		if month, err := strconv.Atoi(monthInput); err == nil && month >= 1 && month <= 12 {
			options.Month = month
		} else {
			fmt.Println("Invalid month, generating calendar for the whole year")
		}
	} else {
		options.Month = 0
	}

	// Get week start day
	fmt.Print("Week starts on (monday, sunday, etc., default: monday): ")
	weekStartInput, _ := reader.ReadString('\n')
	weekStartInput = strings.TrimSpace(weekStartInput)
	if weekStartInput != "" {
		options.FirstDayOfWeek = parseWeekday(weekStartInput)
	}

	// Show calendar week?
	fmt.Print("Show calendar week numbers? (y/n, default: y): ")
	weekInput, _ := reader.ReadString('\n')
	weekInput = strings.TrimSpace(strings.ToLower(weekInput))
	if weekInput == "" || weekInput == "y" || weekInput == "yes" {
		options.ShowCalendarWeek = true
	} else {
		options.ShowCalendarWeek = false
	}

	// Show work week only?
	fmt.Print("Show work week only? (y/n, default: n): ")
	workweekInput, _ := reader.ReadString('\n')
	workweekInput = strings.TrimSpace(strings.ToLower(workweekInput))
	if workweekInput == "y" || workweekInput == "yes" {
		options.ShowWeekends = false
	} else {
		options.ShowWeekends = true
	}

	// Show comments column?
	fmt.Print("Add comments column? (y/n, default: y): ")
	commentsInput, _ := reader.ReadString('\n')
	commentsInput = strings.TrimSpace(strings.ToLower(commentsInput))
	if commentsInput == "" || commentsInput == "y" || commentsInput == "yes" {
		options.ShowComments = true
	} else {
		options.ShowComments = false
	}

	fmt.Println() // Add a blank line before calendar output
}

func printCalendar(options CalendarOptions) {
	if options.Month == 0 {
		// Generate calendar for the whole year
		for m := 1; m <= 12; m++ {
			options.Month = m
			fmt.Println(generateMonthCalendar(options))
			fmt.Println() // Add a blank line between months
		}
	} else {
		// Generate calendar for the specific month
		fmt.Println(generateMonthCalendar(options))
	}
}

func parseWeekday(day string) time.Weekday {
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

func padRight(s string, width int) string {
	if len(s) < width {
		return s + strings.Repeat(" ", width-len(s))
	}
	return s
}

func separatorCell(width int, justify string) string {
	if width <= 0 {
		return ""
	}
	switch strings.ToLower(justify) {
	case "center":
		if width <= 3 {
			return ":-:"
		}
		return ":" + strings.Repeat("-", width-2) + ":"
	case "right":
		return strings.Repeat("-", width-1) + ":"
	default: // left
		return ":" + strings.Repeat("-", width-1)
	}
}

func generateMonthCalendar(options CalendarOptions) string {
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
		if options.ShowCalendarWeek && strings.ToLower(options.Justify) == "center" {
			columnHeaders = append(columnHeaders, "CW ")
			columnWidths = append(columnWidths, len("CW "))
		} else {
			columnHeaders = append(columnHeaders, "CW")
			columnWidths = append(columnWidths, len("CW"))
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
		sb.WriteString(" " + padRight(h, w) + " |")
	}
	sb.WriteString("\n")

	// separator row
	sb.WriteString("|")
	for _, w := range columnWidths {
		sep := separatorCell(w, options.Justify)
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
			cells = append(cells, fmt.Sprintf("%d", w))
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
			sb.WriteString(" " + padRight(cell, w) + " |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
