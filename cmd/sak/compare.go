package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	homeInfo, homeErr := os.Stat(homePath)

	if os.IsNotExist(err) {
		return fmt.Errorf("path does not exist in current directory: %s", currentPath)
	}
	if os.IsNotExist(homeErr) {
		return fmt.Errorf("path does not exist in home directory: %s", homePath)
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

	fmt.Printf("Opening diff between:\n")
	fmt.Printf("  Current: %s\n", currentFile)
	fmt.Printf("  Home:    %s\n", homeFile)

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

func compareDirectories(currentDir, homeDir string) error {
	// Collect all files from both directories
	currentFiles := make(map[string]bool)
	homeFiles := make(map[string]bool)

	// Walk current directory
	err := walkWithSymlinks(currentDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(currentDir, path)
			if err != nil {
				return nil
			}
			currentFiles[relPath] = true
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk current directory: %v", err)
	}

	// Walk home directory
	err = walkWithSymlinks(homeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue walking even if there's an error
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(homeDir, path)
			if err != nil {
				return nil
			}
			homeFiles[relPath] = true
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk home directory: %v", err)
	}

	// Find common files
	var commonFiles []string
	for file := range currentFiles {
		if homeFiles[file] {
			commonFiles = append(commonFiles, file)
		}
	}

	if len(commonFiles) == 0 {
		fmt.Println("No common files found between the directories.")
		return nil
	}

	fmt.Printf("\nFound %d common files to compare:\n", len(commonFiles))
	for _, file := range commonFiles {
		fmt.Printf("  %s\n", file)
	}

	fmt.Print("\nDo you want to compare all files? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %v", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))
	if response != "y" && response != "yes" {
		fmt.Println("Comparison cancelled.")
		return nil
	}

	// Compare each common file
	for i, file := range commonFiles {
		currentFile := filepath.Join(currentDir, file)
		homeFile := filepath.Join(homeDir, file)

		fmt.Printf("\n[%d/%d] Comparing: %s\n", i+1, len(commonFiles), file)

		cmd := exec.Command("code", "--diff", currentFile, homeFile)
		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to open diff for %s: %v\n", file, err)
			continue
		}

		if i < len(commonFiles)-1 {
			fmt.Print("Press Enter to continue to next file...")
			reader.ReadString('\n')
		}
	}

	fmt.Printf("\nCompleted comparing %d files.\n", len(commonFiles))
	return nil
}
