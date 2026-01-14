package types

import "time"

// PackageInfo represents Homebrew package information.
type PackageInfo struct {
	Name      string
	Version   string
	URL       string
	Type      string // "formula", "cask", or "unknown"
	Installed bool
	Failed    bool // true if failed to get package info
}

// Record represents a single work record.
type Record struct {
	Date     time.Time
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Normal   bool // if false, use fixed duration (9h, 10-19) instead
}

// WorktimeSummary represents work time statistics for a single period.
type WorktimeSummary struct {
	Period  string        // "day", "week", "month", "quarter", "year"
	Label   string        // "Day", "Week", "Month", etc.
	Average time.Duration // average work duration
	Count   int           // number of work days
	Error   error         // calculation error
}

// WorktimeComparison represents statistics with comparison data.
type WorktimeComparison struct {
	Current  WorktimeSummary // current period data
	Previous WorktimeSummary // previous period data (optional)
}

// BuildInfo contains all build-time information.
type BuildInfo struct {
	Version      string `json:"version"`
	GitCommit    string `json:"gitCommit"`
	GitBranch    string `json:"gitBranch,omitempty"`
	GitTag       string `json:"gitTag,omitempty"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	GOOS         string `json:"goos"`
	GOARCH       string `json:"goarch"`
}
