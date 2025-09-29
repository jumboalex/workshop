package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Create a temporary test file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, World!"

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	got, err := ReadFile(testFile)
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	if got != content {
		t.Errorf("ReadFile() = %q, want %q", got, content)
	}
}

func TestReadFile_NonExistent(t *testing.T) {
	_, err := ReadFile("nonexistent.txt")
	if err == nil {
		t.Error("ReadFile() expected error for non-existent file, got nil")
	}
}

func TestReadCSVFile(t *testing.T) {
	// Create a temporary CSV file
	tmpDir := t.TempDir()
	csvFile := filepath.Join(tmpDir, "test.csv")
	csvContent := "Name,Age\nAlice,30\nBob,25\n"

	err := os.WriteFile(csvFile, []byte(csvContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test CSV file: %v", err)
	}

	lines, err := ReadCSVFile(csvFile)
	if err != nil {
		t.Errorf("ReadCSVFile() error = %v", err)
	}

	expectedLines := 3
	if len(lines) != expectedLines {
		t.Errorf("ReadCSVFile() returned %d lines, want %d", len(lines), expectedLines)
	}

	if lines[0][0] != "Name" || lines[0][1] != "Age" {
		t.Errorf("ReadCSVFile() header = %v, want [Name Age]", lines[0])
	}
}

func TestReadCSVFile_NonExistent(t *testing.T) {
	_, err := ReadCSVFile("nonexistent.csv")
	if err == nil {
		t.Error("ReadCSVFile() expected error for non-existent file, got nil")
	}
}
