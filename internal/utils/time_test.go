package utils

import (
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	type args struct {
		d time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "zero duration",
			args: args{d: 0},
			want: " 0h  0m",
		},
		{
			name: "1 hour 30 minutes",
			args: args{d: time.Hour + 30*time.Minute},
			want: " 1h 30m",
		},
		{
			name: "23 hours 59 minutes",
			args: args{d: 23*time.Hour + 59*time.Minute},
			want: "23h 59m",
		},
		{
			name: "-1 hour 30 minutes",
			args: args{d: -1*time.Hour - 30*time.Minute},
			want: "-1h -30m",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.args.d); got != tt.want {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	input := time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC)
	want := time.Date(2025, 7, 15, 23, 59, 59, 0, time.UTC)
	got := EndOfDay(input)
	if !got.Equal(want) {
		t.Errorf("EndOfDay() = %v, want %v", got, want)
	}
}

func TestStartOfDay(t *testing.T) {
	input := time.Date(2025, 7, 15, 14, 30, 0, 0, time.UTC)
	want := time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC)
	got := StartOfDay(input)
	if !got.Equal(want) {
		t.Errorf("StartOfDay() = %v, want %v", got, want)
	}
}

func TestWeekRange(t *testing.T) {
	tests := []struct {
		name               string
		input              time.Time
		weekStartsOnMonday bool
		wantStart          time.Time
		wantEnd            time.Time
	}{
		{
			name:               "Wednesday with Monday start",
			input:              time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC), // Wednesday
			weekStartsOnMonday: true,
			wantStart:          time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), // Monday
			wantEnd:            time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC),
		},
		{
			name:               "Wednesday with Sunday start",
			input:              time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC), // Wednesday
			weekStartsOnMonday: false,
			wantStart:          time.Date(2025, 7, 13, 0, 0, 0, 0, time.UTC), // Sunday
			wantEnd:            time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC),
		},
		{
			name:               "Monday with Monday start",
			input:              time.Date(2025, 7, 14, 10, 0, 0, 0, time.UTC), // Monday
			weekStartsOnMonday: true,
			wantStart:          time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), // Same Monday
			wantEnd:            time.Date(2025, 7, 14, 23, 59, 59, 0, time.UTC),
		},
		{
			name:               "Sunday with Monday start",
			input:              time.Date(2025, 7, 20, 10, 0, 0, 0, time.UTC), // Sunday
			weekStartsOnMonday: true,
			wantStart:          time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC), // Previous Monday
			wantEnd:            time.Date(2025, 7, 20, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := WeekRange(tt.input, tt.weekStartsOnMonday)
			if !got.Start.Equal(tt.wantStart) {
				t.Errorf("WeekRange().Start = %v, want %v", got.Start, tt.wantStart)
			}
			if !got.End.Equal(tt.wantEnd) {
				t.Errorf("WeekRange().End = %v, want %v", got.End, tt.wantEnd)
			}
		})
	}
}

func TestLastWeekRange(t *testing.T) {
	// Wednesday 2025-07-16
	input := time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC)
	got := LastWeekRange(input, true)

	// Previous week: Monday 2025-07-07 to Sunday 2025-07-13
	wantStart := time.Date(2025, 7, 7, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 7, 13, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("LastWeekRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("LastWeekRange().End = %v, want %v", got.End, wantEnd)
	}
}

func TestMonthRange(t *testing.T) {
	input := time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC)
	got := MonthRange(input)

	wantStart := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("MonthRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("MonthRange().End = %v, want %v", got.End, wantEnd)
	}
}

func TestLastMonthRange(t *testing.T) {
	input := time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC)
	got := LastMonthRange(input)

	wantStart := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("LastMonthRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("LastMonthRange().End = %v, want %v", got.End, wantEnd)
	}
}

func TestQuarterRange(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart time.Time
	}{
		{
			name:      "Q1",
			input:     time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Q2",
			input:     time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, 4, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Q3",
			input:     time.Date(2025, 7, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:      "Q4",
			input:     time.Date(2025, 11, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := QuarterRange(tt.input)
			if !got.Start.Equal(tt.wantStart) {
				t.Errorf("QuarterRange().Start = %v, want %v", got.Start, tt.wantStart)
			}
		})
	}
}

func TestLastQuarterRange(t *testing.T) {
	tests := []struct {
		name      string
		input     time.Time
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "Q2 -> Q1",
			input:     time.Date(2025, 5, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2025, 3, 31, 23, 59, 59, 0, time.UTC),
		},
		{
			name:      "Q1 -> Q4 of previous year",
			input:     time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			wantStart: time.Date(2024, 10, 1, 0, 0, 0, 0, time.UTC),
			wantEnd:   time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LastQuarterRange(tt.input)
			if !got.Start.Equal(tt.wantStart) {
				t.Errorf("LastQuarterRange().Start = %v, want %v", got.Start, tt.wantStart)
			}
			if !got.End.Equal(tt.wantEnd) {
				t.Errorf("LastQuarterRange().End = %v, want %v", got.End, tt.wantEnd)
			}
		})
	}
}

func TestYearRange(t *testing.T) {
	input := time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC)
	got := YearRange(input)

	wantStart := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("YearRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("YearRange().End = %v, want %v", got.End, wantEnd)
	}
}

func TestLastYearRange(t *testing.T) {
	input := time.Date(2025, 7, 16, 14, 0, 0, 0, time.UTC)
	got := LastYearRange(input)

	wantStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("LastYearRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("LastYearRange().End = %v, want %v", got.End, wantEnd)
	}
}

func TestDayRange(t *testing.T) {
	input := time.Date(2025, 7, 16, 14, 30, 0, 0, time.UTC)
	got := DayRange(input)

	wantStart := time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC)
	wantEnd := time.Date(2025, 7, 16, 23, 59, 59, 0, time.UTC)

	if !got.Start.Equal(wantStart) {
		t.Errorf("DayRange().Start = %v, want %v", got.Start, wantStart)
	}
	if !got.End.Equal(wantEnd) {
		t.Errorf("DayRange().End = %v, want %v", got.End, wantEnd)
	}
}
