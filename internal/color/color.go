// Package color provides colored text output functions.
package color

import "github.com/fatih/color"

// Color functions for text formatting
var (
	Gray   = color.New(color.FgHiBlack).SprintFunc()
	Green  = color.New(color.FgGreen).SprintFunc()
	Yellow = color.New(color.FgYellow).SprintFunc()
	Red    = color.New(color.FgRed).SprintFunc()
	Blue   = color.New(color.FgBlue).SprintFunc()
)

// Disable disables colored output.
func Disable() {
	color.NoColor = true
}

// Enable enables colored output.
func Enable() {
	color.NoColor = false
}
