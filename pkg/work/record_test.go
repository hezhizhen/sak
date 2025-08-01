package work

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func Test_parseDate(t *testing.T) {
	type args struct {
		dateStr string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name:    "with workday",
			args:    args{dateStr: "2025-07-16 Wednesday"},
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "without workday",
			args:    args{dateStr: "2025-07-16"},
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "empty string",
			args:    args{dateStr: ""},
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid date format",
			args:    args{dateStr: "invalid-date"},
			want:    time.Time{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.args.dateStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseTimeOnDate(t *testing.T) {
	baseDate := time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		date    time.Time
		timeStr string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "valid time",
			date:    baseDate,
			timeStr: "09:30:00",
			want:    time.Date(2025, 7, 16, 9, 30, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "midnight",
			date:    baseDate,
			timeStr: "00:00:00",
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "end of day",
			date:    baseDate,
			timeStr: "23:59:59",
			want:    time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "invalid format - missing seconds",
			date:    baseDate,
			timeStr: "09:30",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "invalid format - too many parts",
			date:    baseDate,
			timeStr: "09:30:00:00",
			want:    time.Time{},
			wantErr: true,
		},
		{
			name:    "hour 25 normalizes to next day",
			date:    baseDate,
			timeStr: "25:30:00",
			want:    time.Date(2025, 7, 17, 1, 30, 0, 0, time.UTC), // 25:30 = 01:30 next day
			wantErr: false,
		},
		{
			name:    "minute 70 normalizes",
			date:    baseDate,
			timeStr: "09:70:00",
			want:    time.Date(2025, 7, 16, 10, 10, 0, 0, time.UTC), // 09:70 = 10:10
			wantErr: false,
		},
		{
			name:    "second 70 normalizes",
			date:    baseDate,
			timeStr: "09:30:70",
			want:    time.Date(2025, 7, 16, 9, 31, 10, 0, time.UTC), // 09:30:70 = 09:31:10
			wantErr: false,
		},
		{
			name:    "non-numeric hour",
			date:    baseDate,
			timeStr: "xx:30:00",
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTimeOnDate(tt.date, tt.timeStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseTimeOnDate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseTimeOnDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSingleRecord(t *testing.T) {
	tests := []struct {
		name     string
		dateStr  string
		startStr string
		endStr   string
		want     Record
		wantErr  bool
	}{
		{
			name:     "valid record",
			dateStr:  "2025-07-16",
			startStr: "09:00:00",
			endStr:   "17:30:00",
			want: Record{
				Date:     time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				Start:    time.Date(2025, 7, 16, 9, 0, 0, 0, time.UTC),
				End:      time.Date(2025, 7, 16, 17, 30, 0, 0, time.UTC),
				Duration: 8*time.Hour + 30*time.Minute,
			},
			wantErr: false,
		},
		{
			name:     "record with day name",
			dateStr:  "2025-07-16 Wednesday",
			startStr: "08:30:15",
			endStr:   "16:45:30",
			want: Record{
				Date:     time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				Start:    time.Date(2025, 7, 16, 8, 30, 15, 0, time.UTC),
				End:      time.Date(2025, 7, 16, 16, 45, 30, 0, time.UTC),
				Duration: 8*time.Hour + 15*time.Minute + 15*time.Second,
			},
			wantErr: false,
		},
		{
			name:     "overnight work (end time next day)",
			dateStr:  "2025-07-16",
			startStr: "22:00:00",
			endStr:   "06:00:00",
			want: Record{
				Date:     time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
				Start:    time.Date(2025, 7, 16, 22, 0, 0, 0, time.UTC),
				End:      time.Date(2025, 7, 17, 6, 0, 0, 0, time.UTC),
				Duration: 8 * time.Hour,
			},
			wantErr: false,
		},
		{
			name:     "invalid date",
			dateStr:  "invalid-date",
			startStr: "09:00:00",
			endStr:   "17:00:00",
			want:     Record{},
			wantErr:  true,
		},
		{
			name:     "invalid start time",
			dateStr:  "2025-07-16",
			startStr: "invalid-time",
			endStr:   "17:00:00",
			want:     Record{},
			wantErr:  true,
		},
		{
			name:     "invalid end time",
			dateStr:  "2025-07-16",
			startStr: "09:00:00",
			endStr:   "invalid-time",
			want:     Record{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSingleRecord(tt.dateStr, tt.startStr, tt.endStr)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSingleRecord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseSingleRecord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseRecordsFromFile(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	tests := []struct {
		name        string
		csvContent  string
		expectedLen int
		wantErr     bool
	}{
		{
			name: "valid CSV file",
			csvContent: `Date,Start,End
2025-07-16 Wednesday,09:00:00,17:30:00
2025-07-17 Thursday,08:30:00,16:45:00`,
			expectedLen: 2,
			wantErr:     false,
		},
		{
			name: "single record",
			csvContent: `Date,Start,End
2025-07-16,09:00:00,17:00:00`,
			expectedLen: 1,
			wantErr:     false,
		},
		{
			name: "header only",
			csvContent: `Date,Start,End`,
			expectedLen: 0,
			wantErr:     true,
		},
		{
			name: "invalid CSV format - wrong column count",
			csvContent: `Date,Start,End
2025-07-16,09:00:00`,
			expectedLen: 0,
			wantErr:     true,
		},
		{
			name: "invalid date in record",
			csvContent: `Date,Start,End
invalid-date,09:00:00,17:00:00`,
			expectedLen: 0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filename := filepath.Join(tmpDir, fmt.Sprintf("test_%s.csv", tt.name))
			err := os.WriteFile(filename, []byte(tt.csvContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test the function
			got, err := ParseRecordsFromFile(filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseRecordsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) != tt.expectedLen {
				t.Errorf("ParseRecordsFromFile() got %d records, want %d", len(got), tt.expectedLen)
			}

			// Verify first record for valid cases
			if !tt.wantErr && len(got) > 0 {
				if got[0].Date.IsZero() {
					t.Errorf("ParseRecordsFromFile() first record has zero date")
				}
				if got[0].Duration <= 0 {
					t.Errorf("ParseRecordsFromFile() first record has non-positive duration: %v", got[0].Duration)
				}
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		_, err := ParseRecordsFromFile("/non/existent/file.csv")
		if err == nil {
			t.Errorf("ParseRecordsFromFile() should return error for non-existent file")
		}
	})
}

func TestCalculateAverageForRecords(t *testing.T) {
	// Create test records
	records := []Record{
		{
			Date:     time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), // Monday
			Duration: 8 * time.Hour,
		},
		{
			Date:     time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC), // Tuesday
			Duration: 7*time.Hour + 30*time.Minute,
		},
		{
			Date:     time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC), // Wednesday
			Duration: 8*time.Hour + 30*time.Minute,
		},
		{
			Date:     time.Date(2025, 7, 17, 0, 0, 0, 0, time.UTC), // Thursday
			Duration: 9 * time.Hour,
		},
		{
			Date:     time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC), // Sunday (next week)
			Duration: 6 * time.Hour,
		},
	}

	tests := []struct {
		name            string
		records         []Record
		start           time.Time
		end             time.Time
		expectedAverage time.Duration
		expectedCount   int
		wantErr         bool
	}{
		{
			name:            "all records in range",
			records:         records,
			start:           time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			end:             time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
			expectedAverage: 7*time.Hour + 48*time.Minute, // (8+7.5+8.5+9+6)/5 = 7.8 hours
			expectedCount:   5,
			wantErr:         false,
		},
		{
			name:            "partial range",
			records:         records,
			start:           time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			end:             time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			expectedAverage: 8 * time.Hour, // (7.5+8.5)/2 = 8 hours
			expectedCount:   2,
			wantErr:         false,
		},
		{
			name:            "no records in range",
			records:         records,
			start:           time.Date(2025, 8, 1, 0, 0, 0, 0, time.UTC),
			end:             time.Date(2025, 8, 31, 0, 0, 0, 0, time.UTC),
			expectedAverage: 0,
			expectedCount:   0,
			wantErr:         true,
		},
		{
			name:            "single record",
			records:         records[:1],
			start:           time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			end:             time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			expectedAverage: 8 * time.Hour,
			expectedCount:   1,
			wantErr:         false,
		},
		{
			name:            "empty records",
			records:         []Record{},
			start:           time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			end:             time.Date(2025, 7, 20, 0, 0, 0, 0, time.UTC),
			expectedAverage: 0,
			expectedCount:   0,
			wantErr:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Redirect stdout to avoid test output pollution
			// Note: CalculateAverageForRecords prints to stdout, which we can't easily capture in a unit test
			// For a more thorough test, we would need to refactor the function to accept an io.Writer
			gotAverage, gotCount, err := CalculateAverageForRecords(tt.records, tt.start, tt.end)

			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateAverageForRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if gotCount != tt.expectedCount {
					t.Errorf("CalculateAverageForRecords() gotCount = %v, want %v", gotCount, tt.expectedCount)
				}

				if gotAverage != tt.expectedAverage {
					t.Errorf("CalculateAverageForRecords() gotAverage = %v, want %v", gotAverage, tt.expectedAverage)
				}
			}
		})
	}
}
