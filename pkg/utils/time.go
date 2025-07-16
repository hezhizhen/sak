package utils

import (
	"fmt"
	"time"
)

// FormatDuration formats a time.Duration into a string showing hours and minutes.
// It returns a string in the format "Xh Ym" where X is hours and Y
func FormatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%2dh %2dm", hours, minutes)
}
