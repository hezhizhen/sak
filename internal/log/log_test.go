package log

import (
	"bytes"
	"strings"
	"testing"
)

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, WARN, false)

	if logger == nil {
		t.Fatal("NewLogger returned nil")
	}
	if logger.level != WARN {
		t.Errorf("level = %v, want %v", logger.level, WARN)
	}
	if logger.output != buf {
		t.Error("output not set correctly")
	}
	if logger.colors != false {
		t.Error("colors should be false")
	}
}

func TestLogger_log_LevelFiltering(t *testing.T) {
	tests := []struct {
		name       string
		loggerLvl  Level
		msgLvl     Level
		wantOutput bool
	}{
		{"DEBUG logger, DEBUG msg", DEBUG, DEBUG, true},
		{"DEBUG logger, INFO msg", DEBUG, INFO, true},
		{"INFO logger, DEBUG msg", INFO, DEBUG, false},
		{"INFO logger, INFO msg", INFO, INFO, true},
		{"WARN logger, INFO msg", WARN, INFO, false},
		{"WARN logger, WARN msg", WARN, WARN, true},
		{"ERROR logger, WARN msg", ERROR, WARN, false},
		{"ERROR logger, ERROR msg", ERROR, ERROR, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf, tt.loggerLvl, false)
			logger.log(tt.msgLvl, "test message")

			hasOutput := buf.Len() > 0
			if hasOutput != tt.wantOutput {
				t.Errorf("hasOutput = %v, want %v", hasOutput, tt.wantOutput)
			}
		})
	}
}

func TestLogger_LogMethods(t *testing.T) {
	tests := []struct {
		name   string
		method func(*Logger, string, ...interface{})
		level  string
	}{
		{"Debug", (*Logger).Debug, "DEBUG"},
		{"Info", (*Logger).Info, "INFO"},
		{"Warn", (*Logger).Warn, "WARN"},
		{"Error", (*Logger).Error, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf, DEBUG, false)
			tt.method(logger, "hello %s", "world")

			output := buf.String()
			if !strings.Contains(output, "["+tt.level+"]") {
				t.Errorf("output missing level prefix, got: %s", output)
			}
			if !strings.Contains(output, "hello world") {
				t.Errorf("output missing message, got: %s", output)
			}
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := NewLogger(buf, ERROR, false)

	logger.Debug("should not appear")
	if buf.Len() > 0 {
		t.Error("DEBUG should be filtered at ERROR level")
	}

	logger.SetLevel(DEBUG)
	logger.Debug("should appear")
	if buf.Len() == 0 {
		t.Error("DEBUG should appear after SetLevel(DEBUG)")
	}
}

func TestGlobalFunctions(t *testing.T) {
	oldLogger := defaultLogger
	defer func() { defaultLogger = oldLogger }()

	buf := &bytes.Buffer{}
	defaultLogger = NewLogger(buf, DEBUG, false)

	tests := []struct {
		name   string
		fn     func(string, ...interface{})
		level  string
	}{
		{"Debug", Debug, "DEBUG"},
		{"Info", Info, "INFO"},
		{"Warn", Warn, "WARN"},
		{"Error", Error, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			tt.fn("test")
			if !strings.Contains(buf.String(), "["+tt.level+"]") {
				t.Errorf("global %s missing level prefix", tt.name)
			}
		})
	}
}

func TestGlobalSetLevel(t *testing.T) {
	oldLogger := defaultLogger
	defer func() { defaultLogger = oldLogger }()

	buf := &bytes.Buffer{}
	defaultLogger = NewLogger(buf, ERROR, false)

	Debug("should not appear")
	if buf.Len() > 0 {
		t.Error("DEBUG should be filtered")
	}

	SetLevel(DEBUG)
	Debug("should appear")
	if buf.Len() == 0 {
		t.Error("DEBUG should appear after SetLevel")
	}
}

func TestSetColors(t *testing.T) {
	oldLogger := defaultLogger
	defer func() { defaultLogger = oldLogger }()

	buf := &bytes.Buffer{}
	defaultLogger = NewLogger(buf, INFO, true)

	SetColors(false)
	if defaultLogger.colors != false {
		t.Error("colors should be false after SetColors(false)")
	}

	SetColors(true)
	if defaultLogger.colors != true {
		t.Error("colors should be true after SetColors(true)")
	}
}

func TestLogger_ColorsOutput(t *testing.T) {
	tests := []struct {
		name   string
		colors bool
	}{
		{"with colors", true},
		{"without colors", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := NewLogger(buf, INFO, tt.colors)
			logger.Info("test")

			output := buf.String()
			if !strings.Contains(output, "INFO") {
				t.Errorf("output should contain INFO, got: %s", output)
			}
			if !strings.Contains(output, "test") {
				t.Errorf("output should contain message, got: %s", output)
			}
		})
	}
}
