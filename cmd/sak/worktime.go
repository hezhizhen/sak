package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/hezhizhen/sak/internal/log"
	"github.com/hezhizhen/sak/internal/types"
	"github.com/hezhizhen/sak/internal/utils"
	"github.com/hezhizhen/sak/internal/worktime"
	"github.com/spf13/cobra"
)

const (
	worktimeFile = "worktime.csv"
)

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
  sak worktime analyze              # 工时分析
  sak worktime analyze --year 2025  # 仅分析指定年份
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWorktime(includeComparison)
		},
	}

	cmd.Flags().BoolVarP(&includeComparison, "include-comparison", "c", false, "包含历史数据对比")

	cmd.AddCommand(worktimeAnalyzeCmd())

	return cmd
}

func worktimeAnalyzeCmd() *cobra.Command {
	var year int

	cmd := &cobra.Command{
		Use:   "analyze",
		Short: "分析工时模式",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWorktimeAnalyze(year)
		},
	}

	cmd.Flags().IntVar(&year, "year", 0, "要分析的年份（默认：全部记录）")

	return cmd
}

func runWorktimeAnalyze(year int) error {
	if _, err := os.Stat(worktimeFile); os.IsNotExist(err) {
		return fmt.Errorf("%q not found in current directory", worktimeFile)
	}

	allRecords, err := worktime.ParseRecordsFromFile(worktimeFile)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %v", worktimeFile, err)
	}

	// Filter by year if specified
	var records []types.Record
	if year > 0 {
		for _, r := range allRecords {
			if r.Date.Year() == year {
				records = append(records, r)
			}
		}
	} else {
		records = allRecords
	}

	if len(records) == 0 {
		return fmt.Errorf("no records found")
	}

	// Separate normal and leave days
	var normalDays, leaveDays []types.Record
	for _, r := range records {
		if r.Normal {
			normalDays = append(normalDays, r)
		} else {
			leaveDays = append(leaveDays, r)
		}
	}

	firstDate := records[0].Date
	lastDate := records[len(records)-1].Date

	fmt.Println("=== Work Time Analysis ===")
	fmt.Println()
	fmt.Printf("Period: %s ~ %s (%d days, %d leave days excluded)\n",
		firstDate.Format("2006-01-02"),
		lastDate.Format("2006-01-02"),
		len(records),
		len(leaveDays))
	fmt.Println()

	// Use only normal days for all analysis
	if len(normalDays) == 0 {
		return fmt.Errorf("no normal work days found")
	}

	// --- Schedule ---
	fmt.Println("--- Schedule ---")
	var totalArrivalMinutes, totalLeaveMinutes int
	var earliestArrival, latestLeave time.Time
	var earliestArrivalDate, latestLeaveDate time.Time
	for i, r := range normalDays {
		arrMin := r.Start.Hour()*60 + r.Start.Minute()
		leaveMin := r.End.Hour()*60 + r.End.Minute()
		totalArrivalMinutes += arrMin
		totalLeaveMinutes += leaveMin
		if i == 0 || arrMin < earliestArrival.Hour()*60+earliestArrival.Minute() {
			earliestArrival = r.Start
			earliestArrivalDate = r.Date
		}
		if i == 0 || leaveMin > latestLeave.Hour()*60+latestLeave.Minute() {
			latestLeave = r.End
			latestLeaveDate = r.Date
		}
	}
	avgArrMin := totalArrivalMinutes / len(normalDays)
	avgLeaveMin := totalLeaveMinutes / len(normalDays)
	fmt.Printf("Avg Arrival:   %02d:%02d\n", avgArrMin/60, avgArrMin%60)
	fmt.Printf("Avg Leave:     %02d:%02d\n", avgLeaveMin/60, avgLeaveMin%60)
	fmt.Printf("Earliest:      %02d:%02d (%s)\n", earliestArrival.Hour(), earliestArrival.Minute(), earliestArrivalDate.Format("2006-01-02"))
	fmt.Printf("Latest Leave:  %02d:%02d (%s)\n", latestLeave.Hour(), latestLeave.Minute(), latestLeaveDate.Format("2006-01-02"))

	var longestDay types.Record
	for i, r := range normalDays {
		if i == 0 || r.Duration > longestDay.Duration {
			longestDay = r
		}
	}
	fmt.Printf("Longest Day:   %s (%s)\n", utils.FormatDuration(longestDay.Duration), longestDay.Date.Format("2006-01-02"))
	fmt.Println()

	// --- Duration Distribution ---
	fmt.Println("--- Duration Distribution ---")
	buckets := []struct {
		label string
		min   float64
		max   float64
	}{
		{" 9-10h", 9, 10},
		{"10-11h", 10, 11},
		{"11-12h", 11, 12},
		{"  12h+", 12, 100},
	}
	counts := make([]int, len(buckets))
	for _, r := range normalDays {
		hours := r.Duration.Hours()
		for i, b := range buckets {
			if hours >= b.min && hours < b.max {
				counts[i]++
				break
			}
		}
	}
	maxCount := 0
	for _, c := range counts {
		if c > maxCount {
			maxCount = c
		}
	}
	maxBarWidth := 30
	for i, b := range buckets {
		barLen := 0
		if maxCount > 0 {
			barLen = counts[i] * maxBarWidth / maxCount
		}
		pct := float64(counts[i]) * 100 / float64(len(normalDays))
		bar := strings.Repeat("█", barLen)
		fmt.Printf("%s: %-*s %3d (%2.0f%%)\n", b.label, maxBarWidth, bar, counts[i], pct)
	}
	fmt.Println()

	// --- By Weekday ---
	fmt.Println("--- By Weekday ---")
	weekdayTotals := make([]time.Duration, 7)
	weekdayCounts := make([]int, 7)
	for _, r := range normalDays {
		wd := r.Date.Weekday()
		weekdayTotals[wd] += r.Duration
		weekdayCounts[wd]++
	}
	weekdays := []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday, time.Sunday}
	var maxAvg time.Duration
	weekdayAvgs := make([]time.Duration, 7)
	for _, wd := range weekdays {
		if weekdayCounts[wd] > 0 {
			weekdayAvgs[wd] = weekdayTotals[wd] / time.Duration(weekdayCounts[wd])
			if weekdayAvgs[wd] > maxAvg {
				maxAvg = weekdayAvgs[wd]
			}
		}
	}
	for _, wd := range weekdays {
		if weekdayCounts[wd] == 0 {
			continue
		}
		barLen := 0
		if maxAvg > 0 {
			barLen = int(weekdayAvgs[wd]) * 20 / int(maxAvg)
		}
		bar := strings.Repeat("█", barLen)
		fmt.Printf("%s: %-*s %s\n", wd.String()[:3], 20, bar, utils.FormatDuration(weekdayAvgs[wd]))
	}

	return nil
}

func runWorktime(includeComparison bool) error {
	// Check if worktime.csv exists
	if _, err := os.Stat(worktimeFile); os.IsNotExist(err) {
		return fmt.Errorf("%q not found in current directory", worktimeFile)
	}

	// Parse CSV file
	records, err := worktime.ParseRecordsFromFile(worktimeFile)
	if err != nil {
		return fmt.Errorf("failed to parse %q: %v", worktimeFile, err)
	}

	now := time.Now()
	log.Info("Today: %s", now.Format(time.DateOnly))

	// Calculate all current period statistics (day, week, month, quarter, year)
	comparisons := make([]types.WorktimeComparison, 0, 5)

	// Day: today vs Yesterday
	currentDay := calculateDayDuration(records, now)
	dayComparison := types.WorktimeComparison{Current: currentDay}
	if includeComparison {
		yesterday := now.AddDate(0, 0, -1)
		dayComparison.Previous = calculateDayDuration(records, yesterday)
	}
	comparisons = append(comparisons, dayComparison)

	// This Week vs Last Week
	currentWeek := calculateThisWeekAverage(records, now)
	weekComparison := types.WorktimeComparison{Current: currentWeek}
	if includeComparison {
		weekComparison.Previous = calculateLastWeekAverage(records, now)
	}
	comparisons = append(comparisons, weekComparison)

	// This Month vs Last Month
	currentMonth := calculateThisMonthAverage(records, now)
	monthComparison := types.WorktimeComparison{Current: currentMonth}
	if includeComparison {
		monthComparison.Previous = calculateLastMonthAverage(records, now)
	}
	comparisons = append(comparisons, monthComparison)

	// This Quarter vs Last Quarter
	currentQuarter := calculateThisQuarterAverage(records, now)
	quarterComparison := types.WorktimeComparison{Current: currentQuarter}
	if includeComparison {
		quarterComparison.Previous = calculateLastQuarterAverage(records, now)
	}
	comparisons = append(comparisons, quarterComparison)

	// This Year vs Last Year
	currentYear := calculateThisYearAverage(records, now)
	yearComparison := types.WorktimeComparison{Current: currentYear}
	if includeComparison {
		yearComparison.Previous = calculateLastYearAverage(records, now)
	}
	comparisons = append(comparisons, yearComparison)

	// Format and output table
	output := formatWorktimeTable(comparisons, includeComparison)
	fmt.Print(output)

	return nil
}

// calculatePeriodAverage 通用周期计算函数
func calculatePeriodAverage(records []types.Record, period, label string, tr utils.TimeRange) types.WorktimeSummary {
	average, count, err := worktime.CalculateAverageForRecords(records, tr.Start, tr.End)
	return types.WorktimeSummary{
		Period:  period,
		Label:   label,
		Average: average,
		Count:   count,
		Error:   err,
	}
}

// calculateDayDuration 计算指定日期的工作时长
func calculateDayDuration(records []types.Record, targetDate time.Time) types.WorktimeSummary {
	targetYear, targetMonth, targetDay := targetDate.Date()

	for _, record := range records {
		recordYear, recordMonth, recordDay := record.Date.Date()

		if recordYear == targetYear && recordMonth == targetMonth && recordDay == targetDay {
			return types.WorktimeSummary{
				Period:  "day",
				Label:   "Day",
				Average: record.Duration,
				Count:   1,
				Error:   nil,
			}
		}
	}

	return types.WorktimeSummary{
		Period: "day",
		Label:  "Day",
		Error:  fmt.Errorf("no work time data found for date %s", targetDate.Format(time.DateOnly)),
	}
}

// calculateThisWeekAverage 计算本周平均工作时长（周一为起始日）
func calculateThisWeekAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.WeekRange(now, true) // weekStartsOnMonday = true
	return calculatePeriodAverage(records, "week", "Week", tr)
}

// calculateThisMonthAverage 计算本月平均工作时长
func calculateThisMonthAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.MonthRange(now)
	return calculatePeriodAverage(records, "month", "Month", tr)
}

// calculateThisQuarterAverage 计算本季度平均工作时长
func calculateThisQuarterAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.QuarterRange(now)
	return calculatePeriodAverage(records, "quarter", "Quarter", tr)
}

// calculateLastQuarterAverage 计算上季度平均工作时长
func calculateLastQuarterAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.LastQuarterRange(now)
	return calculatePeriodAverage(records, "quarter", "Quarter", tr)
}

// calculateThisYearAverage 计算今年平均工作时长
func calculateThisYearAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.YearRange(now)
	return calculatePeriodAverage(records, "year", "Year", tr)
}

// calculateLastWeekAverage 计算上周平均工作时长（周一为起始日）
func calculateLastWeekAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.LastWeekRange(now, true) // weekStartsOnMonday = true
	return calculatePeriodAverage(records, "week", "Week", tr)
}

// calculateLastMonthAverage 计算上月平均工作时长
func calculateLastMonthAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.LastMonthRange(now)
	return calculatePeriodAverage(records, "month", "Month", tr)
}

// calculateLastYearAverage 计算去年平均工作时长
func calculateLastYearAverage(records []types.Record, now time.Time) types.WorktimeSummary {
	tr := utils.LastYearRange(now)
	return calculatePeriodAverage(records, "year", "Year", tr)
}

// formatWorktimeTable 格式化工作时间统计表格
func formatWorktimeTable(comparisons []types.WorktimeComparison, includeComparison bool) string {
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
			log.Debug("no data for %s: %v", current.Period, current.Error)
			currentStr = "-"
		} else {
			currentStr = utils.FormatDuration(current.Average)
		}

		if includeComparison {
			// Format previous period duration
			var previousStr string
			if comp.Previous.Error != nil {
				log.Debug("no data for previous %s: %v", comp.Previous.Period, comp.Previous.Error)
				previousStr = "-"
			} else {
				previousStr = utils.FormatDuration(comp.Previous.Average)
			}

			fmt.Fprintf(&result, "%-8s %-12s %s\n",
				current.Label,
				currentStr,
				previousStr)
		} else {
			fmt.Fprintf(&result, "%-8s %s\n",
				current.Label,
				currentStr)
		}
	}

	return result.String()
}
