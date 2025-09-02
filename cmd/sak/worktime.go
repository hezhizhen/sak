package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hezhizhen/sak/pkg/log"
	"github.com/hezhizhen/sak/pkg/utils"
	"github.com/hezhizhen/sak/pkg/work"
	"github.com/spf13/cobra"
)

const (
	worktimeFile = "worktime.csv"
)

// WorktimeSummary 表示单个时间范围的工作时间统计结果
type WorktimeSummary struct {
	Period  string        // "day", "week", "month", "quarter", "year"
	Label   string        // "Today", "This Week", "This Month", etc.
	Average time.Duration // 平均工作时长
	Count   int           // 工作日天数
	Error   error         // 计算错误
}

// WorktimeComparison 表示包含对比数据的统计结果
type WorktimeComparison struct {
	Current  WorktimeSummary // 当前期间数据
	Previous WorktimeSummary // 历史对比数据（可选）
}

func worktimeCmd() *cobra.Command {
	var includeComparison bool

	cmd := &cobra.Command{
		Use:   "worktime",
		Short: fmt.Sprintf("Analyze work time data from %q", worktimeFile),
		Long: `显示工作时间统计，默认输出今天、本周、本月、本季度、今年的数据。

Examples:
  sak worktime                      # 显示当前时间范围统计
  sak worktime -c                   # 显示当前统计 + 历史对比
  sak worktime --include-comparison # 同上
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWorktime(includeComparison)
		},
	}

	cmd.Flags().BoolVarP(&includeComparison, "include-comparison", "c", false, "包含历史数据对比")

	return cmd
}

func runWorktime(includeComparison bool) error {
	// Check if worktime.csv exists
	if _, err := os.Stat(worktimeFile); os.IsNotExist(err) {
		return fmt.Errorf("%q not found in current directory", worktimeFile)
	}

	// Parse CSV file
	records, err := work.ParseRecordsFromFile(worktimeFile)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %v", worktimeFile, err)
	}

	now := time.Now()
	log.Info("Today: %s", now.Format(time.DateOnly))

	// Calculate all current period statistics
	var comparisons []WorktimeComparison

	// Day: today vs Yesterday
	currentDay := calculateDayDuration(records, now)
	dayComparison := WorktimeComparison{Current: currentDay}
	if includeComparison {
		yesterday := now.AddDate(0, 0, -1)
		dayComparison.Previous = calculateDayDuration(records, yesterday)
	}
	comparisons = append(comparisons, dayComparison)

	// This Week vs Last Week
	currentWeek := calculateThisWeekAverage(records, now)
	weekComparison := WorktimeComparison{Current: currentWeek}
	if includeComparison {
		weekComparison.Previous = calculateLastWeekAverage(records, now)
	}
	comparisons = append(comparisons, weekComparison)

	// This Month vs Last Month
	currentMonth := calculateThisMonthAverage(records, now)
	monthComparison := WorktimeComparison{Current: currentMonth}
	if includeComparison {
		monthComparison.Previous = calculateLastMonthAverage(records, now)
	}
	comparisons = append(comparisons, monthComparison)

	// This Quarter vs Last Quarter
	currentQuarter := calculateThisQuarterAverage(records, now)
	quarterComparison := WorktimeComparison{Current: currentQuarter}
	if includeComparison {
		quarterComparison.Previous = calculateLastQuarterAverage(records, now)
	}
	comparisons = append(comparisons, quarterComparison)

	// This Year vs Last Year
	currentYear := calculateThisYearAverage(records, now)
	yearComparison := WorktimeComparison{Current: currentYear}
	if includeComparison {
		yearComparison.Previous = calculateLastYearAverage(records, now)
	}
	comparisons = append(comparisons, yearComparison)

	// Format and output table
	output := formatWorktimeTable(comparisons, includeComparison)
	fmt.Print(output)

	return nil
}

// calculateDayDuration 计算指定日期的工作时长
func calculateDayDuration(records []work.Record, targetDate time.Time) WorktimeSummary {
	targetYear, targetMonth, targetDay := targetDate.Date()

	for _, record := range records {
		recordYear, recordMonth, recordDay := record.Date.Date()

		if recordYear == targetYear && recordMonth == targetMonth && recordDay == targetDay {
			return WorktimeSummary{
				Period:  "day",
				Label:   "Day",
				Average: record.Duration,
				Count:   1,
				Error:   nil,
			}
		}
	}

	return WorktimeSummary{
		Period: "day",
		Label:  "Day",
		Error:  fmt.Errorf("no work time data found for date %s", targetDate.Format(time.DateOnly)),
	}
}

// calculateThisWeekAverage 计算本周平均工作时长
func calculateThisWeekAverage(records []work.Record, now time.Time) WorktimeSummary {
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := now
	endOfWeek = time.Date(endOfWeek.Year(), endOfWeek.Month(), endOfWeek.Day(), 23, 59, 59, 0, endOfWeek.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfWeek, endOfWeek)
	return WorktimeSummary{
		Period:  "week",
		Label:   "Week",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateThisMonthAverage 计算本月平均工作时长
func calculateThisMonthAverage(records []work.Record, now time.Time) WorktimeSummary {
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := now
	endOfMonth = time.Date(endOfMonth.Year(), endOfMonth.Month(), endOfMonth.Day(), 23, 59, 59, 0, endOfMonth.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfMonth, endOfMonth)
	return WorktimeSummary{
		Period:  "month",
		Label:   "Month",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateThisQuarterAverage 计算本季度平均工作时长
func calculateThisQuarterAverage(records []work.Record, now time.Time) WorktimeSummary {
	startOfQuarter, endOfQuarter := getQuarterRange(now.Year(), getQuarter(now), now.Location())
	endOfQuarter = minTime(endOfQuarter, now)
	endOfQuarter = time.Date(endOfQuarter.Year(), endOfQuarter.Month(), endOfQuarter.Day(), 23, 59, 59, 0, endOfQuarter.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfQuarter, endOfQuarter)
	return WorktimeSummary{
		Period:  "quarter",
		Label:   "Quarter",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateLastQuarterAverage 计算上季度平均工作时长
func calculateLastQuarterAverage(records []work.Record, now time.Time) WorktimeSummary {
	currentQuarter := getQuarter(now)
	var lastQuarter int
	var year int

	if currentQuarter == 1 {
		lastQuarter = 4
		year = now.Year() - 1
	} else {
		lastQuarter = currentQuarter - 1
		year = now.Year()
	}

	startOfLastQuarter, endOfLastQuarter := getQuarterRange(year, lastQuarter, now.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfLastQuarter, endOfLastQuarter)
	return WorktimeSummary{
		Period:  "quarter",
		Label:   "Quarter",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateThisYearAverage 计算今年平均工作时长
func calculateThisYearAverage(records []work.Record, now time.Time) WorktimeSummary {
	startOfYear := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	endOfYear := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfYear, endOfYear)
	return WorktimeSummary{
		Period:  "year",
		Label:   "Year",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateLastWeekAverage 计算上周平均工作时长
func calculateLastWeekAverage(records []work.Record, now time.Time) WorktimeSummary {
	weekday := int(now.Weekday())
	startOfThisWeek := now.AddDate(0, 0, -weekday)
	startOfLastWeek := startOfThisWeek.AddDate(0, 0, -7)
	startOfLastWeek = time.Date(startOfLastWeek.Year(), startOfLastWeek.Month(), startOfLastWeek.Day(), 0, 0, 0, 0, startOfLastWeek.Location())

	endOfLastWeek := startOfThisWeek.Add(-time.Second)
	endOfLastWeek = time.Date(endOfLastWeek.Year(), endOfLastWeek.Month(), endOfLastWeek.Day(), 23, 59, 59, 0, endOfLastWeek.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfLastWeek, endOfLastWeek)
	return WorktimeSummary{
		Period:  "week",
		Label:   "Week",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateLastMonthAverage 计算上月平均工作时长
func calculateLastMonthAverage(records []work.Record, now time.Time) WorktimeSummary {
	startOfThisMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	startOfLastMonth := startOfThisMonth.AddDate(0, -1, 0)

	endOfLastMonth := startOfThisMonth.Add(-time.Second)

	average, count, err := work.CalculateAverageForRecords(records, startOfLastMonth, endOfLastMonth)
	return WorktimeSummary{
		Period:  "month",
		Label:   "Month",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateLastYearAverage 计算去年平均工作时长
func calculateLastYearAverage(records []work.Record, now time.Time) WorktimeSummary {
	startOfLastYear := time.Date(now.Year()-1, time.January, 1, 0, 0, 0, 0, now.Location())
	endOfLastYear := time.Date(now.Year()-1, time.December, 31, 23, 59, 59, 0, now.Location())

	average, count, err := work.CalculateAverageForRecords(records, startOfLastYear, endOfLastYear)
	return WorktimeSummary{
		Period:  "year",
		Label:   "Year",
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// getQuarter 根据时间获取季度（1-4）
func getQuarter(t time.Time) int {
	month := int(t.Month())
	return (month-1)/3 + 1
}

// getQuarterRange 获取指定年份和季度的起止时间
func getQuarterRange(year, quarter int, loc *time.Location) (time.Time, time.Time) {
	var startMonth, endMonth time.Month
	switch quarter {
	case 1:
		startMonth, endMonth = time.January, time.March
	case 2:
		startMonth, endMonth = time.April, time.June
	case 3:
		startMonth, endMonth = time.July, time.September
	case 4:
		startMonth, endMonth = time.October, time.December
	}

	start := time.Date(year, startMonth, 1, 0, 0, 0, 0, loc)
	end := time.Date(year, endMonth+1, 1, 0, 0, 0, 0, loc).Add(-time.Second)

	return start, end
}

// minTime 返回两个时间中较早的一个
func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

// formatWorktimeTable 格式化工作时间统计表格
func formatWorktimeTable(comparisons []WorktimeComparison, includeComparison bool) string {
	var result strings.Builder

	// Table headers
	if includeComparison {
		result.WriteString("Period   This Period  Last Period\n")
	} else {
		result.WriteString("Period   Duration\n")
	}

	// Table rows
	for _, comp := range comparisons {
		current := comp.Current

		// Format current period duration
		var currentStr string
		if current.Error != nil {
			log.Error("failed to format current period duration: %v", current.Error)
			currentStr = "-"
		} else {
			currentStr = utils.FormatDuration(current.Average)
		}

		if includeComparison {
			// Format previous period duration
			var previousStr string
			if comp.Previous.Error != nil {
				log.Error("failed to format previous period duration: %v", comp.Previous.Error)
				previousStr = "-"
			} else {
				previousStr = utils.FormatDuration(comp.Previous.Average)
			}

			result.WriteString(fmt.Sprintf("%-8s %-12s %s\n",
				current.Label,
				currentStr,
				previousStr))
		} else {
			result.WriteString(fmt.Sprintf("%-8s %s\n",
				current.Label,
				currentStr))
		}
	}

	return result.String()
}
