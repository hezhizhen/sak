package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hezhizhen/sak/internal/version"
	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the sak version information",
		Long: `Show the sak version information

Example - print version:
  sak version
  sak version --json
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(jsonOutput)
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "Output version information in JSON format")

	return cmd
}

func runVersion(jsonOutput bool) error {
	buildInfo := version.GetBuildInfo()

	if jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(buildInfo)
	}

	fmt.Println(buildInfo.Version)
	return nil
}
