package main

import (
	"os"

	"github.com/hezhizhen/sak/pkg/log"
	"github.com/spf13/cobra"
)

var verbose bool

func main() {
	cmd := &cobra.Command{
		Use:   "sak",
		Short: "My tool set",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add global --verbose flag
	cmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Enable verbose output (show debug messages)")

	// Initialize logger before running any command
	cobra.OnInitialize(initLogger)

	cmd.AddCommand(versionCmd())
	cmd.AddCommand(worktimeCmd())
	cmd.AddCommand(compareCmd())

	if err := cmd.Execute(); err != nil {
		log.Error("%v", err)
		os.Exit(1)
	}
}

func initLogger() {
	if verbose {
		log.SetLevel(log.DEBUG)
	} else {
		log.SetLevel(log.INFO)
	}
}
