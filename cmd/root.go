package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/andre-a-alves/mdcal/pkg/calendar"
	"github.com/andre-a-alves/mdcal/pkg/interactive"
	"github.com/andre-a-alves/mdcal/pkg/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
	rootCmd.PersistentFlags().StringP("justify", "j", "left", "Cell justification: left, center, or right")

	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		// Handle version flag
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Printf("mdcal v%v\n", getVersion())
			os.Exit(0)
		}

		// Get flag values
		weekStart, _ := cmd.Flags().GetString("start")
		noWeekNo, _ := cmd.Flags().GetBool("no-week-no")
		workweek, _ := cmd.Flags().GetBool("workweek")
		noComment, _ := cmd.Flags().GetBool("no-comment")
		justify, _ := cmd.Flags().GetString("justify")

		// Initialize options with defaults
		options := calendar.NewOptions()
		options.FirstDayOfWeek = utils.ParseWeekday(weekStart)
		options.ShowCalendarWeek = !noWeekNo
		options.ShowWeekends = !workweek
		options.ShowComments = !noComment
		options.Justify = justify

		// Check if any args or flags were provided (excluding help flag)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if f.Name == "help" {
				return
			}
		})
		noFlagsOrArgs := cmd.Flags().NFlag() == 0 && len(args) == 0

		if noFlagsOrArgs {
			// Run in interactive mode
			interactive.RunInteractiveMode(&options)
		} else {
			// Handle unnamed arguments (year, month, end year, end month)
			if len(args) > 0 {
				if year, err := strconv.Atoi(args[0]); err == nil {
					options.Year = year
				} else {
					fmt.Println("Invalid year, using current year")
				}
			}

			if len(args) > 1 {
				if month, err := strconv.Atoi(args[1]); err == nil && month >= 1 && month <= 12 {
					options.Month = &month
				} else {
					fmt.Println("Invalid month, generating calendar for the whole year")
				}
			}

			// Handle range format: year month endMonth or year month endYear endMonth
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
					options.EndYear = &endYear
					options.EndMonth = &endMonth
				} else {
					if errYear != nil {
						fmt.Println("Invalid end year, ignoring range")
					}
					if errMonth != nil || endMonth < 1 || endMonth > 12 {
						fmt.Println("Invalid end month, ignoring range")
					}
					options.EndYear = nil
					options.EndMonth = nil
				}
			}
		}

		// Generate and print calendar
		fmt.Print(calendar.PrintCalendar(options))
	}
}
