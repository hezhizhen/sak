package work

import (
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
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.Local),
			wantErr: false,
		},
		{
			name:    "without workday",
			args:    args{dateStr: "2025-07-16"},
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.Local),
			wantErr: false,
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

func Test_hasLeave(t *testing.T) {
	// Create test times
	date := time.Date(2025, 8, 15, 0, 0, 0, 0, time.UTC)

	type args struct {
		start time.Time
		end   time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal work day: 10-19",
			args: args{
				start: date.Add(time.Hour * 10),
				end:   date.Add(time.Hour * 19),
			},
			want: false,
		},
		{
			name: "normal work day: 8-17",
			args: args{
				start: date.Add(time.Hour * 8),
				end:   date.Add(time.Hour * 17),
			},
			want: false,
		},
		{
			name: "normal work day: 11-20",
			args: args{
				start: date.Add(time.Hour * 11),
				end:   date.Add(time.Hour * 20),
			},
			want: false,
		},
		{
			name: "leave day: 10-15",
			args: args{
				start: date.Add(time.Hour * 10),
				end:   date.Add(time.Hour * 15),
			},
			want: true,
		},
		{
			name: "leave day: 6-15",
			args: args{
				start: date.Add(time.Hour * 6),
				end:   date.Add(time.Hour * 15),
			},
			want: true,
		},
		{
			name: "leave day: 14-23",
			args: args{
				start: date.Add(time.Hour * 14),
				end:   date.Add(time.Hour * 23),
			},
			want: true,
		},
		{
			name: "leave day: 9-16",
			args: args{
				start: date.Add(time.Hour * 9),
				end:   date.Add(time.Hour * 16),
			},
			want: true,
		},
		{
			name: "leave day: 11-19",
			args: args{
				start: date.Add(time.Hour * 11),
				end:   date.Add(time.Hour * 19),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasLeave(tt.args.start, tt.args.end); got != tt.want {
				t.Errorf("hasLeave() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseSingleRecord(t *testing.T) {
	type args struct {
		dateStr  string
		startStr string
		endStr   string
	}
	tests := []struct {
		name    string
		args    args
		want    Record
		wantErr bool
	}{
		{
			name: "normal work day",
			args: args{
				dateStr:  "2025-08-15 Friday",
				startStr: "09:41:42",
				endStr:   "19:41:43",
			},
			want: Record{
				Date:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.Local),
				Start:    time.Date(2025, 8, 15, 9, 41, 42, 0, time.Local),
				End:      time.Date(2025, 8, 15, 19, 41, 43, 0, time.Local),
				Duration: 10*time.Hour + 1*time.Second,
				Normal:   true,
			},
			wantErr: false,
		},
		{
			name: "leave day",
			args: args{
				dateStr:  "2025-08-15 Friday",
				startStr: "14:30:00",
				endStr:   "20:35:00",
			},
			want: Record{
				Date:     time.Date(2025, 8, 15, 0, 0, 0, 0, time.Local),
				Start:    time.Date(2025, 8, 15, 14, 30, 0, 0, time.Local),
				End:      time.Date(2025, 8, 15, 20, 35, 0, 0, time.Local),
				Duration: 6*time.Hour + 5*time.Minute,
				Normal:   false, // Leave day
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSingleRecord(tt.args.dateStr, tt.args.startStr, tt.args.endStr)
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

func Test_CalculateAverageForRecords(t *testing.T) {
	// Create test records
	normalDay := Record{
		Date:     time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
		Start:    time.Date(2023, 7, 31, 10, 0, 0, 0, time.UTC),
		End:      time.Date(2023, 7, 31, 19, 0, 0, 0, time.UTC),
		Duration: 9 * time.Hour,
		Normal:   true,
	}
	leaveDay := Record{
		Date:     time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
		Start:    time.Date(2023, 8, 1, 10, 0, 0, 0, time.UTC),
		End:      time.Date(2023, 8, 1, 15, 0, 0, 0, time.UTC), // 5 hours actual
		Duration: 5 * time.Hour,
		Normal:   false, // Leave day
	}

	type args struct {
		records []Record
		start   time.Time
		end     time.Time
	}
	tests := []struct {
		name        string
		args        args
		wantAverage time.Duration
		wantCount   int
		wantErr     bool
	}{
		{
			name: "mixed normal and leave days",
			args: args{
				records: []Record{normalDay, leaveDay},
				start:   time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
				end:     time.Date(2023, 8, 1, 0, 0, 0, 0, time.UTC),
			},
			wantAverage: 9 * time.Hour, // (9h + 9h) / 2 = 9h (leave day counted as 9h)
			wantCount:   2,
			wantErr:     false,
		},
		{
			name: "only normal days",
			args: args{
				records: []Record{normalDay},
				start:   time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
				end:     time.Date(2023, 7, 31, 0, 0, 0, 0, time.UTC),
			},
			wantAverage: 9 * time.Hour,
			wantCount:   1,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAverage, gotCount, err := CalculateAverageForRecords(tt.args.records, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Errorf("CalculateAverageForRecords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAverage != tt.wantAverage {
				t.Errorf("CalculateAverageForRecords() gotAverage = %v, want %v", gotAverage, tt.wantAverage)
			}
			if gotCount != tt.wantCount {
				t.Errorf("CalculateAverageForRecords() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
