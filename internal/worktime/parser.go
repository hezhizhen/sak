package worktime

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hezhizhen/sak/internal/types"
)

// ParseRecordsFromFile parses a CSV file containing work records.
// The CSV file syntax is expected to be:
// Date,Start,End
func ParseRecordsFromFile(filename string) ([]types.Record, error) {
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

	var records []types.Record
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

// parseSingleRecord parses a single record from the CSV file.
func parseSingleRecord(dateStr, startStr, endStr string) (types.Record, error) {
	// Parse date (format: "2025-07-16 Wednesday")
	dateParts := strings.Fields(dateStr)
	if len(dateParts) < 1 {
		return types.Record{}, fmt.Errorf("invalid date format: %s", dateStr)
	}

	date, err := time.ParseInLocation("2006-01-02", dateParts[0], time.Local)
	if err != nil {
		return types.Record{}, fmt.Errorf("failed to parse date %s: %v", dateParts[0], err)
	}

	// Parse start time
	startTime, err := parseTimeOnDate(date, startStr)
	if err != nil {
		return types.Record{}, fmt.Errorf("failed to parse start time %s: %v", startStr, err)
	}

	// Parse end time
	endTime, err := parseTimeOnDate(date, endStr)
	if err != nil {
		return types.Record{}, fmt.Errorf("failed to parse end time %s: %v", endStr, err)
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

	return types.Record{
		Date:     date,
		Start:    startTime,
		End:      endTime,
		Duration: duration,
		Normal:   !isLeaveDay, // Normal is true when there's no leave
	}, nil
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
