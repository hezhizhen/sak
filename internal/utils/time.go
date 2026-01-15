package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a time.Duration into a human-readable string.
// It returns a string in the format "Xh Ym" where X is hours and Y is minutes.
// For negative durations, both hours and minutes will be negative.
// Examples:
//   - FormatDuration(1*time.Hour + 30*time.Minute) returns " 1h 30m"
//   - FormatDuration(-1*time.Hour - 30*time.Minute) returns "-1h -30m"
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%2dh %2dm", hours, minutes)
}

// TimeRange represents a time range with start and end times.
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// EndOfDay returns the time at end of day (23:59:59) for the given time.
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
}

// StartOfDay returns the time at start of day (00:00:00) for the given time.
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// WeekRange returns the time range for the week containing the given time.
// If weekStartsOnMonday is true, the week starts on Monday; otherwise, Sunday.
// The end time is set to the end of the given day.
func WeekRange(t time.Time, weekStartsOnMonday bool) TimeRange {
	weekday := int(t.Weekday())
	if weekStartsOnMonday {
		// Adjust: Monday=0, Tuesday=1, ..., Sunday=6
		weekday = (weekday + 6) % 7
	}
	start := StartOfDay(t.AddDate(0, 0, -weekday))
	end := EndOfDay(t)
	return TimeRange{Start: start, End: end}
}

// LastWeekRange returns the time range for the previous week.
func LastWeekRange(t time.Time, weekStartsOnMonday bool) TimeRange {
	weekday := int(t.Weekday())
	if weekStartsOnMonday {
		weekday = (weekday + 6) % 7
	}
	startOfThisWeek := StartOfDay(t.AddDate(0, 0, -weekday))
	start := startOfThisWeek.AddDate(0, 0, -7)
	end := EndOfDay(startOfThisWeek.Add(-time.Second))
	return TimeRange{Start: start, End: end}
}

// MonthRange returns the time range for the month containing the given time.
// The end time is set to the end of the given day.
func MonthRange(t time.Time) TimeRange {
	start := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	end := EndOfDay(t)
	return TimeRange{Start: start, End: end}
}

// LastMonthRange returns the time range for the previous month.
func LastMonthRange(t time.Time) TimeRange {
	startOfThisMonth := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
	start := startOfThisMonth.AddDate(0, -1, 0)
	lastDayOfLastMonth := startOfThisMonth.Add(-time.Second)
	end := EndOfDay(lastDayOfLastMonth)
	return TimeRange{Start: start, End: end}
}

// QuarterRange returns the time range for the quarter containing the given time.
// The end time is set to the end of the given day.
func QuarterRange(t time.Time) TimeRange {
	quarter := (int(t.Month())-1)/3 + 1
	startMonth := time.Month((quarter-1)*3 + 1)
	start := time.Date(t.Year(), startMonth, 1, 0, 0, 0, 0, t.Location())
	end := EndOfDay(t)
	return TimeRange{Start: start, End: end}
}

// LastQuarterRange returns the time range for the previous quarter.
func LastQuarterRange(t time.Time) TimeRange {
	currentQuarter := (int(t.Month())-1)/3 + 1
	year := t.Year()
	lastQuarter := currentQuarter - 1
	if lastQuarter < 1 {
		lastQuarter = 4
		year--
	}
	startMonth := time.Month((lastQuarter-1)*3 + 1)
	endMonth := startMonth + 2

	start := time.Date(year, startMonth, 1, 0, 0, 0, 0, t.Location())
	// Last day of the quarter
	endDate := time.Date(year, endMonth+1, 1, 0, 0, 0, 0, t.Location()).Add(-time.Second)
	end := EndOfDay(endDate)
	return TimeRange{Start: start, End: end}
}

// YearRange returns the time range for the year containing the given time.
// The end time is set to the end of the given day.
func YearRange(t time.Time) TimeRange {
	start := time.Date(t.Year(), time.January, 1, 0, 0, 0, 0, t.Location())
	end := EndOfDay(t)
	return TimeRange{Start: start, End: end}
}

// LastYearRange returns the time range for the previous year.
func LastYearRange(t time.Time) TimeRange {
	start := time.Date(t.Year()-1, time.January, 1, 0, 0, 0, 0, t.Location())
	end := time.Date(t.Year()-1, time.December, 31, 23, 59, 59, 0, t.Location())
	return TimeRange{Start: start, End: end}
}

// DayRange returns the time range for the day containing the given time.
func DayRange(t time.Time) TimeRange {
	return TimeRange{Start: StartOfDay(t), End: EndOfDay(t)}
}
