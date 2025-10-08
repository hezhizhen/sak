// Package version provides versioning information for the sak application.
package version

import (
	"runtime"
)

var (
	// Version is the current version of the sak.
	Version = "dev"
	// GitCommit is the git sha1
	GitCommit = ""
	// GitTreeState is the state of the git tree
	GitTreeState = ""
	// GitBranch is the git branch
	GitBranch = ""
	// GitTag is the git tag
	GitTag = ""
	// BuildDate is when the binary was built
	BuildDate = ""
	// GoVersion is the Go version used to build
	GoVersion = ""
)

// BuildInfo contains all build-time information
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

// GetBuildInfo returns comprehensive build and runtime information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:      Version,
		GitCommit:    GitCommit,
		GitBranch:    GitBranch,
		GitTag:       GitTag,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion: func() string {
			// use runtime version as default
			if GoVersion != "" {
				return GoVersion
			}
			return runtime.Version()
		}(),
		GOOS:   runtime.GOOS,
		GOARCH: runtime.GOARCH,
	}
}
