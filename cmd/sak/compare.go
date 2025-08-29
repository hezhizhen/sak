package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func compareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <file_path>",
		Short: "Compare a file between current directory and home directory",
		Long: `Compare a file between current directory and home directory using VS Code diff.

Example:
  sak compare .bashrc              # Compare current/.bashrc with ~/.bashrc
  sak compare config/app.conf      # Compare current/config/app.conf with ~/config/app.conf
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(args[0])
		},
	}

	return cmd
}

func runCompare(filePath string) error {
	// Get current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}

	// Get home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %v", err)
	}

	// Construct full paths
	currentFile := filepath.Join(currentDir, filePath)
	homeFile := filepath.Join(homeDir, filePath)

	// Check if current directory file exists
	if _, err := os.Stat(currentFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist in current directory: %s", currentFile)
	}

	// Check if home directory file exists
	if _, err := os.Stat(homeFile); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist in home directory: %s", homeFile)
	}

	// Check if code command is available
	if _, err := exec.LookPath("code"); err != nil {
		return fmt.Errorf("VS Code 'code' command not found. Please make sure VS Code is installed and the 'code' command is available in PATH")
	}

	// Execute code --diff command
	cmd := exec.Command("code", "--diff", currentFile, homeFile)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute VS Code diff: %v", err)
	}

	fmt.Printf("Opening diff between:\n")
	fmt.Printf("  Current: %s\n", currentFile)
	fmt.Printf("  Home:    %s\n", homeFile)

	return nil
}
