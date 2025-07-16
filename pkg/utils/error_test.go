package utils

import (
	"fmt"
	"testing"
)

func TestCheckError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name     string
		args     args
		panicked bool
	}{
		{
			name: "nil error",
			args: args{
				err: nil,
			},
			panicked: false,
		},
		{
			name: "non-nil error",
			args: args{
				err: fmt.Errorf("test error"),
			},
			panicked: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.panicked {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("CheckError did not panic for error: %v", tt.args.err)
					}
				}()
			} else {
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("CheckError panicked for nil error")
					}
				}()
			}

			CheckError(tt.args.err)
		})
	}
}
