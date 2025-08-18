package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hezhizhen/sak/pkg/version"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	var jsonOutput bool
	var shortOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the sak version information",
		Long: `Show the sak version information

Example - print version:
  sak version
  sak version --json
  sak version --short
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(jsonOutput, shortOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information in JSON format")
	cmd.Flags().BoolVar(&shortOutput, "short", false, "Output only the version number")

	return cmd
}

func runVersion(jsonOutput, shortOutput bool) error {
	buildInfo := version.GetBuildInfo()

	if shortOutput {
		fmt.Println(buildInfo.Version)
		return nil
	}

	if jsonOutput {
		// Add runtime information for JSON output
		jsonData := buildInfo

		// Add executable path if available
		if execPath, err := os.Executable(); err == nil {
			jsonData.ExecutablePath = execPath
		}

		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(jsonData)
	}

	// Standard formatted output
	items := [][]string{
		{"Version", buildInfo.Version},
	}

	// Build information section
	if buildInfo.BuildDate != "" {
		items = append(items, []string{"Build Date", buildInfo.BuildDate})
	}
	// Git information section
	if buildInfo.GitCommit != "" {
		items = append(items, []string{"Git Commit", buildInfo.GitCommit})
	}
	if buildInfo.GitBranch != "" {
		items = append(items, []string{"Git Branch", buildInfo.GitBranch})
	}
	if buildInfo.GitTag != "" {
		items = append(items, []string{"Git Tag", buildInfo.GitTag})
	}
	if buildInfo.GitTreeState != "" {
		items = append(items, []string{"Git Tree State", buildInfo.GitTreeState})
	}

	// Go information section
	items = append(items, []string{"Go Version (build)", buildInfo.GoVersion})
	items = append(items, []string{"Go Version (runtime)", buildInfo.GoRuntime})
	items = append(items, []string{"Platform", fmt.Sprintf("%s/%s", buildInfo.GOOS, buildInfo.GOARCH)})
	items = append(items, []string{"CPU Count", fmt.Sprintf("%d", buildInfo.NumCPU)})

	// Runtime information
	if execPath, err := os.Executable(); err == nil {
		items = append(items, []string{"Executable Path", execPath})
		if absPath, err := filepath.Abs(execPath); err == nil && absPath != execPath {
			items = append(items, []string{"Absolute Path", absPath})
		}
	}

	// Find the maximum label width for alignment
	maxWidth := 0
	for _, item := range items {
		if len(item[0]) > maxWidth {
			maxWidth = len(item[0])
		}
	}

	// Print formatted output
	for _, item := range items {
		fmt.Printf("%-*s: %s\n", maxWidth, item[0], item[1])
	}

	return nil
}
