package main

import (
	"github.com/hezhizhen/sak/pkg/utils"
	"github.com/spf13/cobra"
)

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

	cmd.AddCommand(versionCmd())
	cmd.AddCommand(worktimeCmd())
	cmd.AddCommand(diaryCmd())

	utils.CheckError(cmd.Execute())
}
