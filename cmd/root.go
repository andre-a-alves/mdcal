package cmd

import (
	"fmt"
	calendar2 "github.com/andre-a-alves/mdcal/cmd/calendar"
	"github.com/andre-a-alves/mdcal/cmd/interactive"
	"github.com/andre-a-alves/mdcal/cmd/utils"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mdcal",
	Short: "mdcal generates a markdown calendar.",
	Long:  "A customized markdown calendar generator that can either be run interactively or with option flags.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Setup the root command
	rootCmd.Use = "mdcal [year] [month] [endMonth|endYear endMonth]"
	rootCmd.Short = "Generate a markdown calendar"
	rootCmd.Long = `mdcal generates a markdown calendar for the specified year and month or range of months.
Examples:
  mdcal 2025 3        - Generate calendar for March 2025
  mdcal 2025 3 5      - Generate calendar for March through May 2025
  mdcal 2025 12 2026 1 - Generate calendar for December 2025 through January 2026

If no arguments are provided, it runs in interactive mode.`

	// Define flags
	rootCmd.PersistentFlags().StringP("start", "s", "monday", "First day of the week (monday/mon)")
	rootCmd.PersistentFlags().BoolP("no-week-no", "w", false, "Leave week numbers off the calendar")
	rootCmd.PersistentFlags().BoolP("workweek", "W", false, "Leave weekends off the calendar")
	rootCmd.PersistentFlags().BoolP("no-comment", "c", false, "Leave the comments column off")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print version information")
	rootCmd.PersistentFlags().BoolP("short", "S", false, "Display calendar with short day names (Mon, Tue, etc.)")
	rootCmd.PersistentFlags().StringP("justify", "j", "left", "Cell justification: left, center, or right")

	// handleVersionFlag checks if the version flag is set and prints the version if it is
	handleVersionFlag := func(cmd *cobra.Command) bool {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("mdcal v%v\n", getVersion())
			return true
		}
		return false
	}

	// initOptionsFromFlags initializes calendar options from command-line flags
	initOptionsFromFlags := func(cmd *cobra.Command) calendar2.Options {
		weekStart, _ := cmd.Flags().GetString("start")
		noWeekNo, _ := cmd.Flags().GetBool("no-week-no")
		workweek, _ := cmd.Flags().GetBool("workweek")
		noComment, _ := cmd.Flags().GetBool("no-comment")
		shortDayNames, _ := cmd.Flags().GetBool("short")
		justify, _ := cmd.Flags().GetString("justify")

		options := calendar2.NewOptions()
		options.FirstDayOfWeek = utils.ParseWeekday(weekStart)
		options.ShowCalendarWeek = !noWeekNo
		options.ShowWeekends = !workweek
		options.ShowComments = !noComment
		options.UseShortDayNames = shortDayNames
		options.Justify = justify

		return options
	}

	// shouldRunInteractively determines if the program should run in interactive mode
	shouldRunInteractively := func(cmd *cobra.Command, args []string) bool {
		return cmd.Flags().NFlag() == 0 && len(args) == 0
	}

	// processYearArg processes the year argument if present
	processYearArg := func(args []string, options *calendar2.Options) {
		if len(args) > 0 {
			if year, err := strconv.Atoi(args[0]); err == nil {
				if year < 1 || year > 9999 {
					fmt.Println("Year must be between 1 and 9999, using current year")
				} else {
					options.Year = year
				}
			} else {
				fmt.Println("Invalid year, using current year")
			}
		}
	}

	// processMonthArg processes the month argument if present
	processMonthArg := func(args []string, options *calendar2.Options) {
		if len(args) > 1 {
			if month, err := strconv.Atoi(args[1]); err == nil && month >= 1 && month <= 12 {
				options.Month = &month
			} else {
				fmt.Println("Invalid month, generating calendar for the whole year")
			}
		}
	}

	// processDateRangeArgs processes date range arguments if present
	processDateRangeArgs := func(args []string, options *calendar2.Options) {
		if len(args) == 3 {
			// If we have 3 args, it's year month endMonth (same year)
			if endMonth, err := strconv.Atoi(args[2]); err == nil && endMonth >= 1 && endMonth <= 12 {
				endYear := options.Year
				options.EndYear = &endYear
				options.EndMonth = &endMonth
			} else {
				fmt.Println("Invalid end month, ignoring range")
			}
		} else if len(args) == 4 {
			// If we have 4 args, it's year month endYear endMonth
			endYear, errYear := strconv.Atoi(args[2])
			endMonth, errMonth := strconv.Atoi(args[3])

			if errYear == nil && errMonth == nil && endMonth >= 1 && endMonth <= 12 {
				if endYear < 1 || endYear > 9999 {
					fmt.Println("End year must be between 1 and 9999, ignoring range")
					options.EndYear = nil
					options.EndMonth = nil
				} else {
					options.EndYear = &endYear
					options.EndMonth = &endMonth
				}
			} else {
				if errYear != nil {
					fmt.Println("Invalid end year, ignoring range")
				} else if endYear < 1 || endYear > 9999 {
					fmt.Println("End year must be between 1 and 9999, ignoring range")
				}
				if errMonth != nil || endMonth < 1 || endMonth > 12 {
					fmt.Println("Invalid end month, ignoring range")
				}
				options.EndYear = nil
				options.EndMonth = nil
			}
		}
	}

	// processCommandLineArgs processes all command-line arguments
	processCommandLineArgs := func(args []string, options *calendar2.Options) {
		processYearArg(args, options)
		processMonthArg(args, options)
		processDateRangeArgs(args, options)
	}

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// Handle version flag
		if handleVersionFlag(cmd) {
			os.Exit(0)
		}

		// Initialize options from flags
		options := initOptionsFromFlags(cmd)

		if shouldRunInteractively(cmd, args) {
			// Run in interactive mode
			interactive.RunInteractiveMode(&options)
		} else {
			// Process command-line arguments
			processCommandLineArgs(args, &options)
		}

		// Generate and print calendar
		fmt.Print(calendar2.PrintCalendar(options))
	}
}
