package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hezhizhen/sak/internal/log"
	"github.com/spf13/cobra"
)

func compareCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compare <path>",
		Short: "Compare files or directories between current directory and home directory",
		Long: `Compare files or directories between current directory and home directory using VS Code diff.

Examples:
  sak compare .bashrc              # Compare current/.bashrc with ~/.bashrc
  sak compare config/app.conf      # Compare current/config/app.conf with ~/config/app.conf
  sak compare config/              # Compare all files in current/config/ with ~/config/
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCompare(args[0])
		},
	}

	return cmd
}

func runCompare(targetPath string) error {
	// Check if targetPath is absolute
	if filepath.IsAbs(targetPath) {
		return fmt.Errorf("path must be relative, but got absolute path: %s", targetPath)
	}

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
	currentPath := filepath.Join(currentDir, targetPath)
	homePath := filepath.Join(homeDir, targetPath)

	currentInfo, err := os.Stat(currentPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist in current directory: %s", currentPath)
		}
		return fmt.Errorf("failed to stat current path: %v", err)
	}

	homeInfo, err := os.Stat(homePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("path does not exist in home directory: %s", homePath)
		}
		return fmt.Errorf("failed to stat home path: %v", err)
	}

	// Check if code command is available
	if _, err := exec.LookPath("code"); err != nil {
		return fmt.Errorf("VS Code 'code' command not found. Please make sure VS Code is installed and the 'code' command is available in PATH")
	}

	// Handle directory comparison
	if currentInfo.IsDir() && homeInfo.IsDir() {
		return compareDirectories(currentPath, homePath)
	}

	// Handle file comparison
	if !currentInfo.IsDir() && !homeInfo.IsDir() {
		return compareFiles(currentPath, homePath)
	}

	return fmt.Errorf("path type mismatch: one is directory, other is file")
}

func compareFiles(currentFile, homeFile string) error {
	// Execute code --diff command
	cmd := exec.Command("code", "--diff", currentFile, homeFile)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to execute VS Code diff: %v", err)
	}

	// Note: We intentionally don't wait for VS Code to close.
	// VS Code is a GUI application that runs independently.
	// The process will be reaped by the OS when it terminates.

	log.Info("Opening diff between:")
	log.Info("  Current: %s", currentFile)
	log.Info("  Home:    %s", homeFile)

	return nil
}

// walkWithSymlinks is a custom walk function that follows symlinks
func walkWithSymlinks(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return walkFn(path, info, err)
		}

		// If it's a symlink, resolve it and get the actual file info
		if info.Mode()&os.ModeSymlink != 0 {
			resolvedPath, err := filepath.EvalSymlinks(path)
			if err != nil {
				// If we can't resolve the symlink, continue with original info
				return walkFn(path, info, err)
			}

			// Get info of the resolved path
			resolvedInfo, err := os.Stat(resolvedPath)
			if err != nil {
				// If we can't stat the resolved path, continue with original info
				return walkFn(path, info, err)
			}

			// If the resolved path is a directory, we need to walk it
			if resolvedInfo.IsDir() {
				// First call walkFn for the symlink directory itself
				if err := walkFn(path, resolvedInfo, nil); err != nil {
					return err
				}

				// Then walk the contents of the resolved directory
				return filepath.Walk(resolvedPath, func(subPath string, subInfo os.FileInfo, subErr error) error {
					if subErr != nil {
						return walkFn(subPath, subInfo, subErr)
					}

					// Skip the root directory as we already processed it
					if subPath == resolvedPath {
						return nil
					}

					// Calculate relative path from the resolved directory
					relPath, err := filepath.Rel(resolvedPath, subPath)
					if err != nil {
						return err
					}

					// Create the path as if it were under the original symlink
					symlinkSubPath := filepath.Join(path, relPath)
					return walkFn(symlinkSubPath, subInfo, nil)
				})
			} else {
				// It's a symlink to a file, use the resolved info
				return walkFn(path, resolvedInfo, nil)
			}
		}

		// Not a symlink, process normally
		return walkFn(path, info, err)
	})
}

// collectFiles walks a directory and returns a map of relative file paths.
func collectFiles(dir string) (map[string]bool, error) {
	files := make(map[string]bool)
	err := walkWithSymlinks(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			files[relPath] = true
		}
		return nil
	})
	return files, err
}

// findCommonFiles returns files that exist in both maps.
func findCommonFiles(a, b map[string]bool) []string {
	var common []string
	for file := range a {
		if b[file] {
			common = append(common, file)
		}
	}
	return common
}

// promptUserConfirmation asks user for yes/no confirmation.
func promptUserConfirmation(prompt string) (bool, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, fmt.Errorf("failed to read input: %v", err)
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes", nil
}

func compareDirectories(currentDir, homeDir string) error {
	// Collect files from both directories
	currentFiles, err := collectFiles(currentDir)
	if err != nil {
		return fmt.Errorf("failed to walk current directory: %v", err)
	}
	log.Debug("Current directory: %s", currentDir)
	for file := range currentFiles {
		log.Debug("  %s", file)
	}

	homeFiles, err := collectFiles(homeDir)
	if err != nil {
		return fmt.Errorf("failed to walk home directory: %v", err)
	}
	log.Debug("Home directory: %s", homeDir)
	for file := range homeFiles {
		log.Debug("  %s", file)
	}

	// Find common files
	commonFiles := findCommonFiles(currentFiles, homeFiles)
	if len(commonFiles) == 0 {
		log.Info("No common files found between the directories.")
		return nil
	}

	log.Info("Found %d common files to compare:", len(commonFiles))
	for _, file := range commonFiles {
		log.Debug("  %s", file)
	}

	// Ask for user confirmation
	confirmed, err := promptUserConfirmation("\nDo you want to compare all files? (y/n): ")
	if err != nil {
		return err
	}
	if !confirmed {
		log.Info("Comparison cancelled.")
		return nil
	}

	// Compare each common file
	for i, file := range commonFiles {
		currentFile := filepath.Join(currentDir, file)
		homeFile := filepath.Join(homeDir, file)

		log.Info("[%d/%d] Comparing: %s", i+1, len(commonFiles), file)

		cmd := exec.Command("code", "--wait", "--diff", currentFile, homeFile)
		if err := cmd.Run(); err != nil {
			// Log error but continue with remaining files
			log.Error("Failed to run diff for %s: %v", file, err)
			continue
		}
	}

	log.Info("Completed comparing %d files.", len(commonFiles))
	return nil
}
