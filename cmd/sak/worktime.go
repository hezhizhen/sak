package main

import (
	"fmt"
	"os"
	"time"

	"github.com/hezhizhen/sak/pkg/utils"
	"github.com/hezhizhen/sak/pkg/work"
	"github.com/spf13/cobra"
)

func worktimeCmd() *cobra.Command {
	var today, thisWeek, thisMonth, lastWeek, lastMonth, thisYear, lastYear, all bool
	var pastDays int

	cmd := &cobra.Command{
		Use:   "worktime",
		Short: "Analyze work time data from worktime.csv",
		Long: `Analyze work time data from worktime.csv file in the current directory.

Examples:
  sak worktime --today        # Show today's work duration
  sak worktime --this-week    # Show this week's average work duration
  sak worktime --this-month   # Show this month's average work duration
  sak worktime --last-week    # Show last week's average work duration
  sak worktime --last-month   # Show last month's average work duration
  sak worktime --this-year    # Show this year's average work duration
  sak worktime --last-year    # Show last year's average work duration
  sak worktime --all          # Show all time average work duration
  sak worktime --past-days 7  # Show average for past 7 days
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Check if past-days was explicitly set
			pastDaysSet := cmd.Flags().Changed("past-days")
			
			// Ensure only one flag is set
			flags := []bool{today, thisWeek, thisMonth, lastWeek, lastMonth, thisYear, lastYear, all}
			flagCount := 0
			for _, flag := range flags {
				if flag {
					flagCount++
				}
			}

			// Add pastDays to flag count if it's specified
			if pastDaysSet {
				flagCount++
			}

			if flagCount == 0 {
				return fmt.Errorf("please specify one of: --today, --this-week, --this-month, --last-week, --last-month, --this-year, --last-year, --all, or --past-days")
			}

			if flagCount > 1 {
				return fmt.Errorf("please specify only one flag at a time")
			}

			// Validate pastDays value if set
			if pastDaysSet && pastDays <= 0 {
				return fmt.Errorf("past-days must be a positive number")
			}

			return runWorktime(today, thisWeek, thisMonth, lastWeek, lastMonth, thisYear, lastYear, all, pastDays, pastDaysSet)
		},
	}

	cmd.Flags().BoolVar(&today, "today", false, "Show today's work duration")
	cmd.Flags().BoolVar(&thisWeek, "this-week", false, "Show this week's average work duration")
	cmd.Flags().BoolVar(&thisMonth, "this-month", false, "Show this month's average work duration")
	cmd.Flags().BoolVar(&lastWeek, "last-week", false, "Show last week's average work duration")
	cmd.Flags().BoolVar(&lastMonth, "last-month", false, "Show last month's average work duration")
	cmd.Flags().BoolVar(&thisYear, "this-year", false, "Show this year's average work duration")
	cmd.Flags().BoolVar(&lastYear, "last-year", false, "Show last year's average work duration")
	cmd.Flags().BoolVar(&all, "all", false, "Show all time average work duration")
	cmd.Flags().IntVar(&pastDays, "past-days", 0, "Show average for past N days")

	return cmd
}

func runWorktime(today, thisWeek, thisMonth, lastWeek, lastMonth, thisYear, lastYear, all bool, pastDays int, pastDaysSet bool) error {
	// Check if worktime.csv exists
	if _, err := os.Stat("worktime.csv"); os.IsNotExist(err) {
		return fmt.Errorf("worktime.csv not found in current directory")
	}

	// Parse CSV file
	records, err := work.ParseRecordsFromFile("worktime.csv")
	if err != nil {
		return fmt.Errorf("failed to parse worktime.csv: %v", err)
	}

	now := time.Now()

	switch {
	// TODO: support more cases
	case today:
		return showTodayDuration(records, now)
	case thisWeek:
		return showThisWeekAverage(records, now)
	case thisMonth:
		return showThisMonthAverage(records, now)
	case lastWeek:
		return showLastWeekAverage(records, now)
	case lastMonth:
		return showLastMonthAverage(records, now)
	case thisYear:
		return showThisYearAverage(records, now)
	case lastYear:
		return showLastYearAverage(records, now)
	case all:
		return showAllTimeAverage(records)
	case pastDaysSet:
		return showPastDaysAverage(records, now, pastDays)
	}

	return nil
}

func showTodayDuration(records []work.Record, now time.Time) error {
	// Get today's date components
	nowYear, nowMonth, nowDay := now.Date()

	for _, record := range records {
		// Get record's date components
		recordYear, recordMonth, recordDay := record.Date.Date()

		// Compare date components directly
		if recordYear == nowYear && recordMonth == nowMonth && recordDay == nowDay {
			fmt.Printf("%s: %s\n", now.Format(time.DateOnly), utils.FormatDuration(record.Duration))
			return nil
		}
	}

	return fmt.Errorf("no work time data found for today")
}

func showThisWeekAverage(records []work.Record, now time.Time) error {
	// Find the start of this week (Sunday)
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := now
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfWeek, endOfWeek)
	if err != nil {
		return err
	}

	fmt.Printf("This week average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showThisMonthAverage(records []work.Record, now time.Time) error {
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := now
	endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day(), 23, 59, 59, 0, endOfMonth.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfMonth, endOfMonth)
	if err != nil {
		return err
	}

	fmt.Printf("This month average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showLastWeekAverage(records []work.Record, now time.Time) error {
	// Find the start of this week (Sunday), then go back one week
	weekday := int(now.Weekday())
	startOfThisWeek := now.AddDate(0, 0, -weekday)
	startOfLastWeek := startOfThisWeek.AddDate(0, 0, -7)
	startOfLastWeek = time.Date(startOfLastWeek.Year(), startOfLastWeek.Month(), startOfLastWeek.Day(), 0, 0, 0, 0, startOfLastWeek.Location())

	endOfLastWeek := startOfThisWeek.Add(-time.Second)
	endOfLastWeek = time.Date(endOfLastWeek.Year(), endOfLastWeek.Month(), endOfLastWeek.Day(), 23, 59, 59, 0, endOfLastWeek.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfLastWeek, endOfLastWeek)
	if err != nil {
		return err
	}

	fmt.Printf("Last week average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showLastMonthAverage(records []work.Record, now time.Time) error {
	// Get first day of current month, then go back to first day of last month
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)

	// End of last month is one second before start of this month
	endOfLastMonth := startOfThisMonth.Add(-time.Second)

	average, count, err := work.CalculateAverageForRecords(records, startOfLastMonth, endOfLastMonth)
	if err != nil {
		return err
	}

	fmt.Printf("Last month average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showThisYearAverage(records []work.Record, now time.Time) error {
	startOfYear := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	endOfYear := now
	endOfYear = time.Date(endOfYear.Year(), endOfYear.Month(), endOfYear.Day(), 23, 59, 59, 0, endOfYear.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfYear, endOfYear)
	if err != nil {
		return err
	}

	fmt.Printf("This year average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showLastYearAverage(records []work.Record, now time.Time) error {
	startOfLastYear := time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, now.Location())
	endOfLastYear := time.Date(now.Year()-1, 12, 31, 23, 59, 59, 0, now.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfLastYear, endOfLastYear)
	if err != nil {
		return err
	}

	fmt.Printf("Last year average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showAllTimeAverage(records []work.Record) error {
	if len(records) == 0 {
		return fmt.Errorf("no work time data found")
	}

	// Find the earliest and latest dates
	earliest := records[0].Date
	latest := records[0].Date
	for _, record := range records {
		if record.Date.Before(earliest) {
			earliest = record.Date
		}
		if record.Date.After(latest) {
			latest = record.Date
		}
	}

	// Set time bounds for the full range
	startTime := time.Date(earliest.Year(), earliest.Month(), earliest.Day(), 0, 0, 0, 0, earliest.Location())
	endTime := time.Date(latest.Year(), latest.Month(), latest.Day(), 23, 59, 59, 0, latest.Location())

	average, count, err := work.CalculateAverageForRecords(records, startTime, endTime)
	if err != nil {
		return err
	}

	fmt.Printf("All time average (%d days): %s\n", count, utils.FormatDuration(average))
	return nil
}

func showPastDaysAverage(records []work.Record, now time.Time, days int) error {
	if days <= 0 {
		return fmt.Errorf("past-days must be a positive number")
	}

	// Calculate the start date (N days ago from today)
	startDate := now.AddDate(0, 0, -days+1) // +1 because we want to include today
	startTime := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	
	// End time is end of today
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	average, count, err := work.CalculateAverageForRecords(records, startTime, endTime)
	if err != nil {
		return err
	}

	fmt.Printf("Past %d days average (%d days): %s\n", days, count, utils.FormatDuration(average))
	return nil
}
