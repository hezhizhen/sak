package main

import (
	"os"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:          "sak [command]",
		SilenceUsage: true,
	}
	cmd.AddCommand(versionCmd())
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}