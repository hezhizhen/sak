package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func ccusageCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ccusage [claude|amp|opencode|codex]",
		Short: "Run ccusage tools",
		Long: `Run npx ccusage tools.

Supported subcommands:
  claude    - Runs npx ccusage@latest
  amp       - Runs npx @ccusage/amp@latest
  opencode  - Runs npx @ccusage/opencode@latest
  codex     - Runs npx @ccusage/codex@latest`,
		ValidArgs: []string{"claude", "amp", "opencode", "codex"},
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			packages := map[string]string{
				"claude":   "ccusage@latest",
				"amp":      "@ccusage/amp@latest",
				"opencode": "@ccusage/opencode@latest",
				"codex":    "@ccusage/codex@latest",
			}
			packageName := packages[args[0]]

			// We use -y to avoid "Need to install the following packages:" prompt
			c := exec.Command("npx", "-y", packageName)

			// Connect standard streams to allow interaction
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			if err := c.Run(); err != nil {
				return fmt.Errorf("failed to run ccusage tool: %w", err)
			}
			return nil
		},
	}

	return cmd
}
