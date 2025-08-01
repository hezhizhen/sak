package diary

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetDiaryPath(t *testing.T) {
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	
	tests := []struct {
		name     string
		date     time.Time
		baseDir  string
		expected string
	}{
		{
			name:     "default base directory",
			date:     date,
			baseDir:  "",
			expected: "2023-12-25.md",
		},
		{
			name:     "custom base directory",
			date:     date,
			baseDir:  "/tmp/diaries",
			expected: "/tmp/diaries/2023-12-25.md",
		},
	}
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GetDiaryPath(test.date, test.baseDir)
			if result != test.expected {
				t.Errorf("expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestCreateOrAppendEntry(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "diary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	date := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	
	// Test creating new diary file
	err = CreateOrAppendEntry(date, "", tmpDir)
	if err != nil {
		t.Fatalf("failed to create diary entry: %v", err)
	}
	
	// Check if file was created
	diaryPath := GetDiaryPath(date, tmpDir)
	if _, err := os.Stat(diaryPath); os.IsNotExist(err) {
		t.Errorf("diary file was not created")
	}
	
	// Test appending entry
	err = CreateOrAppendEntry(date, "Test entry", tmpDir)
	if err != nil {
		t.Fatalf("failed to append diary entry: %v", err)
	}
	
	// Read file content to verify entry was appended
	content, err := os.ReadFile(diaryPath)
	if err != nil {
		t.Fatalf("failed to read diary file: %v", err)
	}
	
	contentStr := string(content)
	if !contains(contentStr, "Test entry") {
		t.Errorf("entry was not appended to diary file")
	}
}

func TestFindDiaryFiles(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "diary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create some diary files
	diaryFiles := []string{
		"2023-12-25.md",
		"2023-12-26.md",
		"2024-01-01.md",
		"not-a-diary.txt",
	}
	
	for _, file := range diaryFiles {
		filePath := filepath.Join(tmpDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}
	
	// Find diary files
	found, err := FindDiaryFiles(tmpDir)
	if err != nil {
		t.Fatalf("failed to find diary files: %v", err)
	}
	
	// Should find 3 diary files (excluding not-a-diary.txt)
	expectedCount := 3
	if len(found) != expectedCount {
		t.Errorf("expected %d diary files, found %d", expectedCount, len(found))
	}
}

func TestFilterDiariesByDate(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "diary_test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)
	
	// Create diary files with different dates
	diaryFiles := []string{
		"2023-12-24.md",
		"2023-12-25.md",
		"2023-12-26.md",
		"2024-01-01.md",
	}
	
	for _, file := range diaryFiles {
		filePath := filepath.Join(tmpDir, file)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("failed to create test file %s: %v", file, err)
		}
	}
	
	// Filter for Christmas week (Dec 24-26)
	start := time.Date(2023, 12, 24, 0, 0, 0, 0, time.UTC)
	end := time.Date(2023, 12, 26, 23, 59, 59, 0, time.UTC)
	
	filtered, err := FilterDiariesByDate(tmpDir, start, end)
	if err != nil {
		t.Fatalf("failed to filter diary files: %v", err)
	}
	
	// Should find 3 files (Dec 24, 25, 26)
	expectedCount := 3
	if len(filtered) != expectedCount {
		t.Errorf("expected %d filtered files, found %d", expectedCount, len(filtered))
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && 
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || 
		 indexOf(s, substr) >= 0))
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}