// Package log provides structured logging functionality with colored output.
package log

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// Level represents the logging level.
type Level int

const (
	// DEBUG level - only shown with --verbose flag
	DEBUG Level = iota
	// INFO level - shown by default
	INFO
	// WARN level - shown by default
	WARN
	// ERROR level - shown by default
	ERROR
)

// String returns the string representation of the log level.
func (l Level) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Color codes for different log levels
const (
	ColorReset = "\033[0m"
	ColorDebug = "\033[90m" // Gray
	ColorInfo  = "\033[32m" // Green
	ColorWarn  = "\033[33m" // Yellow
	ColorError = "\033[31m" // Red
)

// levelColors maps log levels to their corresponding colors
var levelColors = map[Level]string{
	DEBUG: ColorDebug,
	INFO:  ColorInfo,
	WARN:  ColorWarn,
	ERROR: ColorError,
}

// Logger represents a logger instance.
type Logger struct {
	sync.Mutex

	level  Level
	output io.Writer
	colors bool
}

// NewLogger creates a new logger instance.
func NewLogger(output io.Writer, level Level, colors bool) *Logger {
	return &Logger{
		level:  level,
		output: output,
		colors: colors,
	}
}

// log writes a log message if the level is enabled.
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	l.Lock()
	defer l.Unlock()

	message := fmt.Sprintf(format, args...)

	if l.colors {
		color := levelColors[level]
		fmt.Fprintf(l.output, "%s[%s]%s %s\n", color, level.String(), ColorReset, message)
	} else {
		fmt.Fprintf(l.output, "[%s] %s\n", level.String(), message)
	}
}

// Debug logs a debug message.
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info logs an info message.
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn logs a warning message.
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error logs an error message.
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level Level) {
	l.Lock()
	defer l.Unlock()

	l.level = level
}

// Global logger instance
var defaultLogger *Logger

func init() {
	// Initialize with INFO level and colors enabled by default
	defaultLogger = NewLogger(os.Stdout, INFO, true)
}

// SetLevel sets the global logger level.
func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

// Debug logs a debug message using the global logger.
func Debug(format string, args ...interface{}) {
	defaultLogger.Debug(format, args...)
}

// Info logs an info message using the global logger.
func Info(format string, args ...interface{}) {
	defaultLogger.Info(format, args...)
}

// Warn logs a warning message using the global logger.
func Warn(format string, args ...interface{}) {
	defaultLogger.Warn(format, args...)
}

// Error logs an error message using the global logger.
func Error(format string, args ...interface{}) {
	defaultLogger.Error(format, args...)
}

// SetColors enables or disables colored output.
func SetColors(enabled bool) {
	defaultLogger.Lock()
	defer defaultLogger.Unlock()

	defaultLogger.colors = enabled
}
