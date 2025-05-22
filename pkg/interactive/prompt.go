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
			options.Month = &month
		} else {
			fmt.Println("Invalid month, generating calendar for the whole year")
		}
	} else {
		options.Month = nil
	}

	// Get week start day
	fmt.Print("First day of the week (monday/mon, sunday/sun, etc., default: monday): ")
	weekStartInput, _ := reader.ReadString('\n')
	weekStartInput = strings.TrimSpace(weekStartInput)
	if weekStartInput != "" {
		options.FirstDayOfWeek = utils.ParseWeekday(weekStartInput)
	}

	// Leave week numbers off?
	fmt.Print("Leave week numbers off the calendar? (y/n, default: n): ")
	weekInput, _ := reader.ReadString('\n')
	weekInput = strings.TrimSpace(strings.ToLower(weekInput))
	if weekInput == "y" || weekInput == "yes" {
		options.ShowCalendarWeek = false
	} else {
		options.ShowCalendarWeek = true
	}

	// Leave weekends off?
	fmt.Print("Leave weekends off the calendar? (y/n, default: n): ")
	workweekInput, _ := reader.ReadString('\n')
	workweekInput = strings.TrimSpace(strings.ToLower(workweekInput))
	if workweekInput == "y" || workweekInput == "yes" {
		options.ShowWeekends = false
	} else {
		options.ShowWeekends = true
	}

	// Leave comments column off?
	fmt.Print("Leave comments column off? (y/n, default: n): ")
	commentsInput, _ := reader.ReadString('\n')
	commentsInput = strings.TrimSpace(strings.ToLower(commentsInput))
	if commentsInput == "y" || commentsInput == "yes" {
		options.ShowComments = false
	} else {
		options.ShowComments = true
	}

	// Cell justification
	fmt.Print("Cell justification (left, center, right, default: left): ")
	justifyInput, _ := reader.ReadString('\n')
	justifyInput = strings.TrimSpace(strings.ToLower(justifyInput))
	if justifyInput != "" {
		if justifyInput == "left" || justifyInput == "center" || justifyInput == "right" {
			options.Justify = justifyInput
		} else {
			fmt.Println("Invalid justification, using left")
		}
	}

	fmt.Println() // Add a blank line before calendar output
}
