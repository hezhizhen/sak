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
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "without workday",
			args:    args{dateStr: "2025-07-16"},
			want:    time.Date(2025, 7, 16, 0, 0, 0, 0, time.UTC),
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
