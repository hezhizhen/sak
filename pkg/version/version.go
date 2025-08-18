// Package version provides versioning information for the sak application.
package version

import (
	"runtime"
)

var (
	// Version is the current version of the sak.
	Version = "0.0.1"
	// BuildMetadata is the extra build time data
	BuildMetadata = "unreleased"
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

// GetVersion returns the semver string of the version
func GetVersion() string {
	if BuildMetadata == "" {
		return Version
	}
	return Version + "+" + BuildMetadata
}

// BuildInfo contains all build-time information
type BuildInfo struct {
	Version       string
	GitCommit     string
	GitBranch     string
	GitTag        string
	GitTreeState  string
	BuildDate     string
	GoVersion     string
	GoRuntime     string
	GOOS          string
	GOARCH        string
	NumCPU        int
	BuildMetadata string
}

// GetBuildInfo returns comprehensive build and runtime information
func GetBuildInfo() BuildInfo {
	return BuildInfo{
		Version:       GetVersion(),
		GitCommit:     GitCommit,
		GitBranch:     GitBranch,
		GitTag:        GitTag,
		GitTreeState:  GitTreeState,
		BuildDate:     BuildDate,
		GoVersion:     GoVersion,
		GoRuntime:     runtime.Version(),
		GOOS:          runtime.GOOS,
		GOARCH:        runtime.GOARCH,
		NumCPU:        runtime.NumCPU(),
		BuildMetadata: BuildMetadata,
	}
}
