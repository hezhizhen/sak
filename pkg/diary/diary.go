package diary

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Entry represents a diary entry with timestamp and content
type Entry struct {
	Timestamp time.Time
	Content   string
}

// Diary represents a diary file for a specific date
type Diary struct {
	Date    time.Time
	Path    string
	Entries []Entry
}

// DefaultTemplate is the template used for new diary files
const DefaultTemplate = `# Diary - %s

## Morning Thoughts


## Daily Log


## Evening Reflection


---
`

// GetDiaryPath returns the path for a diary file based on date
func GetDiaryPath(date time.Time, baseDir string) string {
	if baseDir == "" {
		baseDir = "."
	}
	filename := date.Format("2006-01-02") + ".md"
	return filepath.Join(baseDir, filename)
}

// CreateOrAppendEntry creates a new diary file if it doesn't exist, or appends an entry to an existing one
func CreateOrAppendEntry(date time.Time, content string, baseDir string) error {
	diaryPath := GetDiaryPath(date, baseDir)
	
	// Check if file exists
	if _, err := os.Stat(diaryPath); os.IsNotExist(err) {
		// Create new file with template
		template := fmt.Sprintf(DefaultTemplate, date.Format("2006-01-02"))
		if err := os.WriteFile(diaryPath, []byte(template), 0644); err != nil {
			return fmt.Errorf("failed to create diary file: %v", err)
		}
	}
	
	// Append entry with timestamp if content is provided
	if content != "" {
		timestamp := time.Now().Format("15:04")
		entry := fmt.Sprintf("\n**%s** - %s\n", timestamp, content)
		
		file, err := os.OpenFile(diaryPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open diary file: %v", err)
		}
		defer file.Close()
		
		if _, err := file.WriteString(entry); err != nil {
			return fmt.Errorf("failed to write entry: %v", err)
		}
	}
	
	return nil
}

// FindDiaryFiles finds all diary files in the specified directory
func FindDiaryFiles(baseDir string) ([]string, error) {
	if baseDir == "" {
		baseDir = "."
	}
	
	pattern := filepath.Join(baseDir, "????-??-??.md")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find diary files: %v", err)
	}
	
	return files, nil
}

// FilterDiariesByDate filters diary files by date range
func FilterDiariesByDate(baseDir string, start, end time.Time) ([]string, error) {
	files, err := FindDiaryFiles(baseDir)
	if err != nil {
		return nil, err
	}
	
	var filtered []string
	for _, file := range files {
		// Extract date from filename
		basename := filepath.Base(file)
		dateStr := strings.TrimSuffix(basename, ".md")
		
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue // Skip invalid date formats
		}
		
		// Check if date is within range
		if (date.Equal(start) || date.After(start)) && (date.Equal(end) || date.Before(end)) {
			filtered = append(filtered, file)
		}
	}
	
	return filtered, nil
}

// GetDiariesForToday returns diary files for today
func GetDiariesForToday(baseDir string) ([]string, error) {
	today := time.Now()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Second)
	
	return FilterDiariesByDate(baseDir, startOfDay, endOfDay)
}

// GetDiariesForThisWeek returns diary files for this week (Sunday to Saturday)
func GetDiariesForThisWeek(baseDir string) ([]string, error) {
	now := time.Now()
	weekday := int(now.Weekday())
	startOfWeek := now.AddDate(0, 0, -weekday)
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.Add(7*24*time.Hour).Add(-time.Second)
	
	return FilterDiariesByDate(baseDir, startOfWeek, endOfWeek)
}

// GetDiariesForThisMonth returns diary files for this month
func GetDiariesForThisMonth(baseDir string) ([]string, error) {
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)
	
	return FilterDiariesByDate(baseDir, startOfMonth, endOfMonth)
}

// OpenWithEditor opens a file with the specified editor
func OpenWithEditor(filepath string, editor string) error {
	if editor == "" {
		editor = os.Getenv("EDITOR")
		if editor == "" {
			editor = "vim" // Default to vim
		}
	}
	
	// Note: This would require exec.Command in a real implementation
	// For now, just return the command that would be executed
	fmt.Printf("Would execute: %s %s\n", editor, filepath)
	return nil
}