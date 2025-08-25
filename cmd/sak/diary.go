package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/hezhizhen/sak/pkg/diary"
	"github.com/spf13/cobra"
)

func diaryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diary",
		Short: "Manage diary entries",
		Long: `Manage diary entries with support for creating, viewing, and editing diary files.

Diary files are stored as markdown files with the format YYYY-MM-DD.md.

Examples:
  sak diary create                    # Create today's diary entry
  sak diary create "Had a great day"  # Create today's entry with content
  sak diary view --today             # View today's diary entries
  sak diary view --this-week         # View this week's diary entries
  sak diary edit                     # Edit today's diary in default editor
  sak diary edit --editor code       # Edit today's diary in VS Code
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(diaryCreateCmd())
	cmd.AddCommand(diaryViewCmd())
	cmd.AddCommand(diaryEditCmd())

	return cmd
}

func diaryCreateCmd() *cobra.Command {
	var baseDir string

	cmd := &cobra.Command{
		Use:   "create [content]",
		Short: "Create a diary entry for today",
		Long: `Create a diary entry for today. If the file doesn't exist, it will be created with a template.
If content is provided, it will be appended with a timestamp.

Examples:
  sak diary create                    # Create today's diary file with template
  sak diary create "Had lunch"        # Add entry with timestamp
`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			content := ""
			if len(args) > 0 {
				content = args[0]
			}

			today := time.Now()
			err := diary.CreateOrAppendEntry(today, content, baseDir)
			if err != nil {
				return fmt.Errorf("failed to create diary entry: %v", err)
			}

			diaryPath := diary.GetDiaryPath(today, baseDir)
			if content == "" {
				fmt.Printf("Created diary file: %s\n", diaryPath)
			} else {
				fmt.Printf("Added entry to diary: %s\n", diaryPath)
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&baseDir, "dir", "", "Directory to store diary files (default: current directory)")

	return cmd
}

func diaryViewCmd() *cobra.Command {
	var baseDir string
	var today, thisWeek, thisMonth bool

	cmd := &cobra.Command{
		Use:   "view",
		Short: "View diary entries",
		Long: `View diary entries filtered by date range.

Examples:
  sak diary view --today        # View today's diary
  sak diary view --this-week    # View this week's diaries
  sak diary view --this-month   # View this month's diaries
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Ensure only one flag is set
			flags := []bool{today, thisWeek, thisMonth}
			flagCount := 0
			for _, flag := range flags {
				if flag {
					flagCount++
				}
			}

			if flagCount == 0 {
				return fmt.Errorf("please specify one of: --today, --this-week, --this-month")
			}

			if flagCount > 1 {
				return fmt.Errorf("please specify only one flag at a time")
			}

			return runDiaryView(today, thisWeek, thisMonth, baseDir)
		},
	}

	cmd.Flags().StringVar(&baseDir, "dir", "", "Directory to search for diary files (default: current directory)")
	cmd.Flags().BoolVar(&today, "today", false, "View today's diary")
	cmd.Flags().BoolVar(&thisWeek, "this-week", false, "View this week's diaries")
	cmd.Flags().BoolVar(&thisMonth, "this-month", false, "View this month's diaries")

	return cmd
}

func diaryEditCmd() *cobra.Command {
	var baseDir, editor string

	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit today's diary entry",
		Long: `Open today's diary entry in an editor. If the file doesn't exist, it will be created first.

Examples:
  sak diary edit                 # Edit in default editor ($EDITOR or vim)
  sak diary edit --editor code   # Edit in VS Code
  sak diary edit --editor vim    # Edit in vim
`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			today := time.Now()
			diaryPath := diary.GetDiaryPath(today, baseDir)

			// Create diary file if it doesn't exist
			if _, err := os.Stat(diaryPath); os.IsNotExist(err) {
				err := diary.CreateOrAppendEntry(today, "", baseDir)
				if err != nil {
					return fmt.Errorf("failed to create diary file: %v", err)
				}
				fmt.Printf("Created diary file: %s\n", diaryPath)
			}

			// Open in editor
			return openInEditor(diaryPath, editor)
		},
	}

	cmd.Flags().StringVar(&baseDir, "dir", "", "Directory to store diary files (default: current directory)")
	cmd.Flags().StringVar(&editor, "editor", "", "Editor to use (default: $EDITOR or vim)")

	return cmd
}

func runDiaryView(today, thisWeek, thisMonth bool, baseDir string) error {
	var files []string
	var err error

	switch {
	case today:
		files, err = diary.GetDiariesForToday(baseDir)
		if err != nil {
			return fmt.Errorf("failed to get today's diaries: %v", err)
		}
		fmt.Println("Today's diary:")
	case thisWeek:
		files, err = diary.GetDiariesForThisWeek(baseDir)
		if err != nil {
			return fmt.Errorf("failed to get this week's diaries: %v", err)
		}
		fmt.Println("This week's diaries:")
	case thisMonth:
		files, err = diary.GetDiariesForThisMonth(baseDir)
		if err != nil {
			return fmt.Errorf("failed to get this month's diaries: %v", err)
		}
		fmt.Println("This month's diaries:")
	}

	if len(files) == 0 {
		fmt.Println("No diary files found.")
		return nil
	}

	// Display files and their content
	for _, file := range files {
		fmt.Printf("\n--- %s ---\n", filepath.Base(file))
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Error reading file: %v\n", err)
			continue
		}
		fmt.Print(string(content))
	}

	return nil
}

func openInEditor(filepath string, editor string) error {
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default to vim
		}
	}

	// Handle common editors
	var cmd *exec.Cmd
	switch {
	case strings.Contains(editor, "code"):
		cmd = exec.Command(editor, filepath)
	case strings.Contains(editor, "vim") || strings.Contains(editor, "nvim"):
		cmd = exec.Command(editor, filepath)
	case strings.Contains(editor, "nano"):
		cmd = exec.Command(editor, filepath)
	default:
		cmd = exec.Command(editor, filepath)
	}

	// Set up command to use current terminal
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Opening %s with %s...\n", filepath, editor)
	return cmd.Run()
}