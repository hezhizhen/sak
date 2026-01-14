package worktime

import (
	"fmt"
	"time"

	"github.com/hezhizhen/sak/internal/log"
	"github.com/hezhizhen/sak/internal/types"
	"github.com/hezhizhen/sak/internal/utils"
)

// CalculateAverageForRecords calculates the average work time for the given records.
// For leave days (Normal = false), it uses 9 hours in the average calculation instead of actual duration.
func CalculateAverageForRecords(records []types.Record, start, end time.Time) (time.Duration, int, error) {
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
				durationForAverage = MinWorkHours * time.Hour
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
