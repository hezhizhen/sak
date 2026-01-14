package color

import (
	"testing"

	fatihcolor "github.com/fatih/color"
)

func TestColorFunctions(t *testing.T) {
	tests := []struct {
		name  string
		fn    func(a ...interface{}) string
		input string
	}{
		{"Gray", Gray, "test"},
		{"Green", Green, "test"},
		{"Yellow", Yellow, "test"},
		{"Red", Red, "test"},
		{"Blue", Blue, "test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fn(tt.input)
			if result == "" {
				t.Errorf("%s() returned empty string", tt.name)
			}
		})
	}
}

func TestDisable(t *testing.T) {
	Enable()
	Disable()
	if !fatihcolor.NoColor {
		t.Error("Disable() should set NoColor to true")
	}
}

func TestEnable(t *testing.T) {
	Disable()
	Enable()
	if fatihcolor.NoColor {
		t.Error("Enable() should set NoColor to false")
	}
}
