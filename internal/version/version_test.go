package version

import (
	"runtime"
	"testing"
)

func TestGetBuildInfo(t *testing.T) {
	tests := []struct {
		name           string
		setupVersion   string
		setupGoVersion string
		wantVersion    string
		wantGoVersion  string
	}{
		{
			name:           "default values",
			setupVersion:   "dev",
			setupGoVersion: "",
			wantVersion:    "dev",
			wantGoVersion:  runtime.Version(),
		},
		{
			name:           "custom GoVersion",
			setupVersion:   "dev",
			setupGoVersion: "go1.21.0",
			wantVersion:    "dev",
			wantGoVersion:  "go1.21.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalVersion := Version
			originalGoVersion := GoVersion
			defer func() {
				Version = originalVersion
				GoVersion = originalGoVersion
			}()

			Version = tt.setupVersion
			GoVersion = tt.setupGoVersion

			info := GetBuildInfo()

			if info.Version != tt.wantVersion {
				t.Errorf("Version = %q, want %q", info.Version, tt.wantVersion)
			}
			if info.GoVersion != tt.wantGoVersion {
				t.Errorf("GoVersion = %q, want %q", info.GoVersion, tt.wantGoVersion)
			}
			if info.GOOS != runtime.GOOS {
				t.Errorf("GOOS = %q, want %q", info.GOOS, runtime.GOOS)
			}
			if info.GOARCH != runtime.GOARCH {
				t.Errorf("GOARCH = %q, want %q", info.GOARCH, runtime.GOARCH)
			}
		})
	}
}
