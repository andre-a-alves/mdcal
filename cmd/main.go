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
}

func main() {
	// Define flags
	firstDayLong := flag.String("firstday", "monday", "First day of the week (monday, sunday, etc.)")
	showCalWeek := flag.Bool("week", true, "Show calendar week numbers")
	showWeekends := flag.Bool("weekends", true, "Show weekend days")
	workweek := flag.Bool("workweek", false, "Show only workdays (Monday-Friday)")
	showComments := flag.Bool("comments", true, "Add a comments column")
	versionFlag := flag.Bool("version", false, "Print version information")

	// Parse flags
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Println("Calendar v0.1.0")
		os.Exit(0)
	}

	// Handle workweek alias
	if *workweek {
		*showWeekends = false
	}

	// Initialize options with defaults
	options := CalendarOptions{
		Year:             time.Now().Year(),
		Month:            0, // 0 means generate for the whole year
		FirstDayOfWeek:   parseWeekday(*firstDayLong),
		ShowCalendarWeek: *showCalWeek,
		ShowWeekends:     *showWeekends,
		ShowComments:     *showComments,
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

	// Get first day of week
	fmt.Print("First day of week (monday, sunday, etc., default: monday): ")
	firstDayInput, _ := reader.ReadString('\n')
	firstDayInput = strings.TrimSpace(firstDayInput)
	if firstDayInput != "" {
		options.FirstDayOfWeek = parseWeekday(firstDayInput)
	}

	// Show calendar week?
	fmt.Print("Show calendar week numbers? (y/n, default: n): ")
	weekInput, _ := reader.ReadString('\n')
	weekInput = strings.TrimSpace(strings.ToLower(weekInput))
	options.ShowCalendarWeek = weekInput == "y" || weekInput == "yes"

	// Show weekends?
	fmt.Print("Show weekends? (y/n, default: y): ")
	weekendsInput, _ := reader.ReadString('\n')
	weekendsInput = strings.TrimSpace(strings.ToLower(weekendsInput))
	if weekendsInput != "" {
		options.ShowWeekends = !(weekendsInput == "n" || weekendsInput == "no")
	}

	// Show comments column?
	fmt.Print("Add comments column? (y/n, default: n): ")
	commentsInput, _ := reader.ReadString('\n')
	commentsInput = strings.TrimSpace(strings.ToLower(commentsInput))
	options.ShowComments = commentsInput == "y" || commentsInput == "yes"

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

func generateMonthCalendar(options CalendarOptions) string {
	var sb strings.Builder

	// header
	month := time.Month(options.Month)
	sb.WriteString(fmt.Sprintf("# %s %d\n\n", month.String(), options.Year))

	// all weekday names Mon=0…Sun=6
	allNames := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	// convert time.Weekday (Mon=1…) to 0-based Monday=0…Sunday=6
	firstIndex := (int(options.FirstDayOfWeek) + 6) % 7
	shifted := append(allNames[firstIndex:], allNames[:firstIndex]...)

	// strip weekends if needed
	var dayNames []string
	for _, d := range shifted {
		if !options.ShowWeekends && (d == "Sat" || d == "Sun") {
			continue
		}
		dayNames = append(dayNames, d)
	}

	// header row
	if options.ShowCalendarWeek {
		sb.WriteString("| Week |")
	}
	for _, d := range dayNames {
		sb.WriteString(fmt.Sprintf(" %s |", d))
	}
	if options.ShowComments {
		sb.WriteString(" Comments |")
	}
	sb.WriteString("\n")

	// separator
	if options.ShowCalendarWeek {
		sb.WriteString("|------|")
	}
	for range dayNames {
		sb.WriteString("-----|")
	}
	if options.ShowComments {
		sb.WriteString("----------|")
	}
	sb.WriteString("\n")

	// map each column back to a time.Weekday
	var weekDays []time.Weekday
	for _, d := range dayNames {
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

	// compute the first day of the month and its “week start”
	firstOfMonth := time.Date(options.Year, month, 1, 0, 0, 0, 0, time.UTC)
	lastOfMonth := time.Date(options.Year, month+1, 0, 0, 0, 0, 0, time.UTC)
	shift := (int(firstOfMonth.Weekday()) - int(options.FirstDayOfWeek) + 7) % 7
	weekStart := firstOfMonth.AddDate(0, 0, -shift)

	// print each week
	for cur := weekStart; !cur.After(lastOfMonth); cur = cur.AddDate(0, 0, 7) {
		if options.ShowCalendarWeek {
			_, w := cur.ISOWeek()
			sb.WriteString(fmt.Sprintf("| %02d  |", w))
		}
		for _, wd := range weekDays {
			// compute this column’s date
			delta := (int(wd) - int(options.FirstDayOfWeek) + 7) % 7
			cd := cur.AddDate(0, 0, delta)
			if cd.Month() == month {
				sb.WriteString(fmt.Sprintf(" %2d  |", cd.Day()))
			} else {
				sb.WriteString("     |")
			}
		}
		if options.ShowComments {
			sb.WriteString("          |")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
