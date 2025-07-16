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
	var today, thisWeek, thisMonth, lastWeek, lastMonth bool

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
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure only one flag is set
			flags := []bool{today, thisWeek, thisMonth, lastWeek, lastMonth}
			flagCount := 0
			for _, flag := range flags {
				if flag {
					flagCount++
				}
			}

			if flagCount == 0 {
				return fmt.Errorf("please specify one of: --today, --this-week, --this-month, --last-week, --last-month")
			}

			if flagCount > 1 {
				return fmt.Errorf("please specify only one flag at a time")
			}

			return runWorktime(today, thisWeek, thisMonth, lastWeek, lastMonth)
		},
	}

	cmd.Flags().BoolVar(&today, "today", false, "Show today's work duration")
	cmd.Flags().BoolVar(&thisWeek, "this-week", false, "Show this week's average work duration")
	cmd.Flags().BoolVar(&thisMonth, "this-month", false, "Show this month's average work duration")
	cmd.Flags().BoolVar(&lastWeek, "last-week", false, "Show last week's average work duration")
	cmd.Flags().BoolVar(&lastMonth, "last-month", false, "Show last month's average work duration")

	return cmd
}

func runWorktime(today, thisWeek, thisMonth, lastWeek, lastMonth bool) error {
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
