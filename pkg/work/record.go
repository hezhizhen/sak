// Package work provides functionality to manage and parse work records.
package work

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hezhizhen/sak/pkg/log"
	"github.com/hezhizhen/sak/pkg/utils"
)

const (
	afternoonStartHour = 12
	earlyEndHour       = 17
	minWorkHours       = 9
)

type Record struct {
	Date     time.Time
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Normal   bool // if it is false, use fixed duration (9h, 10-19) instead
}

// ParseRecordsFromFile parses a CSV file containing work records.
// The CSV file syntax is expected to be:
// Date,Start,End
func ParseRecordsFromFile(filename string) ([]Record, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("CSV file must have at least a header and one data row")
	}

	var records []Record
	for i, row := range rows {
		if i == 0 { // Skip header
			continue
		}

		if len(row) != 3 {
			return nil, fmt.Errorf("invalid CSV format at row %d: expected 3 columns, got %d", i+1, len(row))
		}

		record, err := parseSingleRecord(row[0], row[1], row[2])
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %v", i+1, err)
		}

		records = append(records, record)
	}

	return records, nil
}

// hasLeave determines if there is leave on a given day based on start and end.
// leave cases:
// - start > 12:00
// - end < 17:00
// - duration < 9 hours
func hasLeave(start, end time.Time) bool {
	afternoonStart := time.Date(start.Year(), start.Month(), start.Day(), afternoonStartHour, 0, 0, 0, start.Location())
	earlyEnd := time.Date(end.Year(), end.Month(), end.Day(), earlyEndHour, 0, 0, 0, end.Location())

	return start.After(afternoonStart) ||
		end.Before(earlyEnd) ||
		end.Sub(start).Hours() < minWorkHours
}

// parseSingleRecord parses a single record from the CSV file.
func parseSingleRecord(dateStr, startStr, endStr string) (Record, error) {
	// Parse date (format: "2025-07-16 Wednesday")
	dateParts := strings.Fields(dateStr)
	if len(dateParts) < 1 {
		return Record{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	date, err := time.ParseInLocation("2006-01-02", dateParts[0], time.Local)
	if err != nil {
		return Record{}, fmt.Errorf("failed to parse date %s: %v", dateParts[0], err)
	}

	// Parse start time
	startTime, err := parseTimeOnDate(date, startStr)
	if err != nil {
		return Record{}, fmt.Errorf("failed to parse start time %s: %v", startStr, err)
	}

	// Parse end time
	endTime, err := parseTimeOnDate(date, endStr)
	if err != nil {
		return Record{}, fmt.Errorf("failed to parse end time %s: %v", endStr, err)
	}

	// Calculate duration
	duration := endTime.Sub(startTime)
	if duration < 0 {
		// Handle case where end time is next day
		endTime = endTime.Add(24 * time.Hour)
		duration = endTime.Sub(startTime)
	}

	// Determine if this day has leave
	isLeaveDay := hasLeave(startTime, endTime)

	return Record{
		Date:     date,
		Start:    startTime,
		End:      endTime,
		Duration: duration,
		Normal:   !isLeaveDay, // Normal is true when there's no leave
	}, nil
}

// parseDate parses a date string to a time.Time.
// The string can be one of the following formats:
// - "2006-01-02 Wednesday"
// - "2006-01-02"
func parseDate(dateStr string) (time.Time, error) {
	dateParts := strings.Fields(dateStr)
	if len(dateParts) < 1 {
		return time.Time{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	date, err := time.ParseInLocation("2006-01-02", dateParts[0], time.Local)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to parse date %s: %v", dateParts[0], err)
	}

	return date, nil
}

func parseTimeOnDate(date time.Time, timeStr string) (time.Time, error) {
	timeParts := strings.Split(timeStr, ":")
	if len(timeParts) != 3 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", timeStr)
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid hour: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid minute: %s", timeParts[1])
	}

	second, err := strconv.Atoi(timeParts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid second: %s", timeParts[2])
	}

	return time.Date(date.Year(), date.Month(), date.Day(), hour, minute, second, 0, date.Location()), nil
}

// CalculateAverageForRecords calculates the average work time for the given records.
// For leave days (Normal = false), it uses 9 hours in the average calculation instead of actual duration.
func CalculateAverageForRecords(records []Record, start, end time.Time) (time.Duration, int, error) {
	var totalDuration time.Duration
	count := 0

	for _, record := range records {
		recordDate := time.Date(record.Date.Year(), record.Date.Month(), record.Date.Day(), 0, 0, 0, 0, record.Date.Location())

		if (recordDate.Equal(start) || recordDate.After(start)) && (recordDate.Equal(end) || recordDate.Before(end)) {
			// Display each selected day's work time
			log.Debug("%2d %s: %s", count+1, record.Date.Format("2006-01-02"), utils.FormatDuration(record.Duration))

			// For average calculation: use 9h for leave days, actual duration for normal days
			var durationForAverage time.Duration
			if record.Normal {
				durationForAverage = record.Duration
			} else {
				// Use 9 hours for leave days in average calculation
				durationForAverage = minWorkHours * time.Hour
			}

			totalDuration += durationForAverage
			count++
		}
	}

	if count == 0 {
		return 0, 0, fmt.Errorf("no work time data found for the specified period")
	}

	average := totalDuration / time.Duration(count)
	return average, count, nil
}
