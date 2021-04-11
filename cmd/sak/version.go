package main

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/hezhizhen/sak/pkg/version"

	"github.com/spf13/cobra"
)

func versionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show the sak version information",
		Long: `Show the sak version information

Example - print version:
  sak version
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion()
		},
	}

	return cmd
}

func runVersion() error {
	items := [][]string{
		{"Version", version.GetVersion()},
		{"Go version", runtime.Version()},
	}
	if version.GitCommit != "" {
		items = append(items, []string{"Git commit", version.GitCommit})
	}
	if version.GitTreeState != "" {
		items = append(items, []string{"Git tree state", version.GitTreeState})
	}

	size := 0
	for _, item := range items {
		if length := len(item[0]); length > size {
			size = length
		}
	}
	for _, item := range items {
		fmt.Println(item[0] + ": " + strings.Repeat(" ", size-len(item[0])) + item[1])
	}

	return nil
}
