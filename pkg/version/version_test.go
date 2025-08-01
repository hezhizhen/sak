package version

import (
	"testing"
)

func TestGetVersion(t *testing.T) {
	tests := []struct {
		name          string
		version       string
		buildMetadata string
		expected      string
	}{
		{
			name:          "version with build metadata",
			version:       "1.0.0",
			buildMetadata: "beta.1",
			expected:      "1.0.0+beta.1",
		},
		{
			name:          "version with empty build metadata",
			version:       "1.0.0",
			buildMetadata: "",
			expected:      "1.0.0",
		},
		{
			name:          "version with unreleased build metadata",
			version:       "0.0.1",
			buildMetadata: "unreleased",
			expected:      "0.0.1+unreleased",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalVersion := Version
			originalBuildMetadata := BuildMetadata

			// Set test values
			Version = tt.version
			BuildMetadata = tt.buildMetadata

			// Test the function
			got := GetVersion()

			// Restore original values
			Version = originalVersion
			BuildMetadata = originalBuildMetadata

			if got != tt.expected {
				t.Errorf("GetVersion() = %v, want %v", got, tt.expected)
			}
		})
	}
}