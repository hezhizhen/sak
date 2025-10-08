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
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "zero duration",
			args:    args{d: 0},
			want:    " 0h  0m",
			wantErr: false,
		},
		{
			name:    "1 hour 30 minutes",
			args:    args{d: time.Hour + 30*time.Minute},
			want:    " 1h 30m",
			wantErr: false,
		},
		{
			name:    "23 hours 59 minutes",
			args:    args{d: 23*time.Hour + 59*time.Minute},
			want:    "23h 59m",
			wantErr: false,
		},
		{
			name:    "-1 hour 30 minutes",
			args:    args{d: -1*time.Hour - 30*time.Minute},
			want:    "-1h -30m",
			wantErr: false,
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
