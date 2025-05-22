package interactive

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/andre-a-alves/mdcal/pkg/calendar"
	"github.com/andre-a-alves/mdcal/pkg/utils"
)

// RunInteractiveMode prompts the user for calendar options and updates the provided options
func RunInteractiveMode(options *calendar.Options) {
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
		options.FirstDayOfWeek = utils.ParseWeekday(weekStartInput)
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