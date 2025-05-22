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

// readInput reads a line of input from the reader and trims whitespace
func readInput(reader *bufio.Reader) string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// promptForYear prompts the user for a year and updates the options
func promptForYear(reader *bufio.Reader, options *calendar.Options) {
	fmt.Printf("Enter year (default: %d): ", options.Year)
	yearInput := readInput(reader)

	if yearInput != "" {
		if year, err := strconv.Atoi(yearInput); err == nil {
			options.Year = year
		} else {
			fmt.Println("Invalid year, using current year")
		}
	}
}

// promptForMonth prompts the user for a month and updates the options
func promptForMonth(reader *bufio.Reader, options *calendar.Options) {
	fmt.Print("Enter month (1-12, empty for whole year): ")
	monthInput := readInput(reader)

	if monthInput != "" {
		if month, err := strconv.Atoi(monthInput); err == nil && month >= 1 && month <= 12 {
			options.Month = &month
		} else {
			fmt.Println("Invalid month, generating calendar for the whole year")
		}
	} else {
		options.Month = nil
	}
}

// promptForWeekStartDay prompts the user for the first day of the week and updates the options
func promptForWeekStartDay(reader *bufio.Reader, options *calendar.Options) {
	fmt.Print("First day of the week (monday/mon, sunday/sun, etc., default: monday): ")
	weekStartInput := readInput(reader)

	if weekStartInput != "" {
		options.FirstDayOfWeek = utils.ParseWeekday(weekStartInput)
	}
}

// promptForBooleanOption prompts the user for a yes/no option and returns the result
func promptForBooleanOption(reader *bufio.Reader, prompt string) bool {
	fmt.Print(prompt)
	input := readInput(reader)
	input = strings.ToLower(input)

	return !(input == "y" || input == "yes")
}

// promptForJustification prompts the user for cell justification and updates the options
func promptForJustification(reader *bufio.Reader, options *calendar.Options) {
	fmt.Print("Cell justification (left, center, right, default: left): ")
	justifyInput := readInput(reader)
	justifyInput = strings.ToLower(justifyInput)

	if justifyInput != "" {
		if justifyInput == "left" || justifyInput == "center" || justifyInput == "right" {
			options.Justify = justifyInput
		} else {
			fmt.Println("Invalid justification, using left")
		}
	}
}

// RunInteractiveMode prompts the user for calendar options and updates the provided options
func RunInteractiveMode(options *calendar.Options) {
	reader := bufio.NewReader(os.Stdin)

	// Prompt for each option
	promptForYear(reader, options)
	promptForMonth(reader, options)
	promptForWeekStartDay(reader, options)

	// Boolean options
	options.ShowCalendarWeek = promptForBooleanOption(reader, "Leave week numbers off the calendar? (y/n, default: n): ")
	options.ShowWeekends = promptForBooleanOption(reader, "Leave weekends off the calendar? (y/n, default: n): ")
	options.ShowComments = promptForBooleanOption(reader, "Leave comments column off? (y/n, default: n): ")

	// Justification
	promptForJustification(reader, options)

	fmt.Println() // Add a blank line before calendar output
}
