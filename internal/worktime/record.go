// Package worktime provides functionality to manage and parse work records.
package worktime

import "time"

// Record represents a single work record
type Record struct {
	Date     time.Time
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Normal   bool // if it is false, use fixed duration (9h, 10-19) instead
}

// hasLeave determines if there is leave on a given day based on start and end.
// leave cases:
// - start > 12:00
// - end < 17:00
// - duration < 9 hours
func hasLeave(start, end time.Time) bool {
	afternoonStart := time.Date(start.Year(), start.Month(), start.Day(), AfternoonStartHour, 0, 0, 0, start.Location())
	earlyEnd := time.Date(end.Year(), end.Month(), end.Day(), EarlyEndHour, 0, 0, 0, end.Location())

	return start.After(afternoonStart) ||
		end.Before(earlyEnd) ||
		end.Sub(start).Hours() < MinWorkHours
}
