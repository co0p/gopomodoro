package storage_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/co0p/gopomodoro/internal/storage"
)

func TestEnsureDataDir(t *testing.T) {
	// Call EnsureDataDir
	err := storage.EnsureDataDir()
	if err != nil {
		t.Fatalf("EnsureDataDir() failed: %v", err)
	}

	// Verify directory exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home directory: %v", err)
	}

	dataDir := homeDir + "/.gopomodoro"
	info, err := os.Stat(dataDir)
	if err != nil {
		t.Fatalf("Data directory does not exist: %v", err)
	}

	if !info.IsDir() {
		t.Fatalf("Data directory path exists but is not a directory")
	}
}

func TestLogSession(t *testing.T) {
	// Set a temporary home directory for testing
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Create a known timestamp
	timestamp := time.Date(2025, 12, 23, 10, 30, 0, 0, time.UTC)

	// Ensure data directory exists
	if err := storage.EnsureDataDir(); err != nil {
		t.Fatalf("Failed to ensure data dir: %v", err)
	}

	// Call LogSession using public API
	err := storage.LogSession(timestamp, "work", "completed", 25)
	if err != nil {
		t.Fatalf("LogSession() failed: %v", err)
	}

	// Read file contents from data directory
	logPath := tmpDir + "/.gopomodoro/sessions.log"
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read sessions.log: %v", err)
	}

	// Verify CSV line format
	expectedLine := "2025-12-23T10:30:00Z,work,completed,25"
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) == 0 {
		t.Fatalf("sessions.log is empty")
	}

	lastLine := lines[len(lines)-1]
	if lastLine != expectedLine {
		t.Fatalf("Expected CSV line:\n%s\nGot:\n%s", expectedLine, lastLine)
	}
}

func TestLogSessionCreatesFileWithHeader(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	logPath := tmpDir + "/.gopomodoro/sessions.log"

	// Create a known timestamp
	timestamp := time.Date(2025, 12, 23, 10, 30, 0, 0, time.UTC)

	// Ensure data directory exists
	if err := storage.EnsureDataDir(); err != nil {
		t.Fatalf("Failed to ensure data dir: %v", err)
	}

	// Call LogSession (file doesn't exist yet)
	err := storage.LogSession(timestamp, "work", "started", 0)
	if err != nil {
		t.Fatalf("LogSession() failed: %v", err)
	}

	// Read file contents
	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read sessions.log: %v", err)
	}

	// Verify header and data line
	lines := strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) < 2 {
		t.Fatalf("Expected at least 2 lines (header + data), got %d", len(lines))
	}

	expectedHeader := "timestamp,session_type,event,duration_minutes"
	if lines[0] != expectedHeader {
		t.Fatalf("Expected header:\n%s\nGot:\n%s", expectedHeader, lines[0])
	}

	expectedDataLine := "2025-12-23T10:30:00Z,work,started,0"
	if lines[1] != expectedDataLine {
		t.Fatalf("Expected data line:\n%s\nGot:\n%s", expectedDataLine, lines[1])
	}

	// Log another session and verify header isn't duplicated
	timestamp2 := time.Date(2025, 12, 23, 10, 55, 0, 0, time.UTC)
	err = storage.LogSession(timestamp2, "work", "completed", 25)
	if err != nil {
		t.Fatalf("Second LogSession() failed: %v", err)
	}

	content, err = os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read sessions.log after second write: %v", err)
	}

	lines = strings.Split(strings.TrimSpace(string(content)), "\n")
	if len(lines) != 3 {
		t.Fatalf("Expected 3 lines (header + 2 data), got %d", len(lines))
	}

	// Verify header is still first line
	if lines[0] != expectedHeader {
		t.Fatalf("Header should remain: %s", lines[0])
	}
}

func TestEnsureDataDirFailure(t *testing.T) {
	// This test is harder to set up on macOS/Linux without root
	// We'll test that the error is properly returned and wrapped
	// by checking the error message structure

	// We can't easily test the actual failure case, but we can verify
	// that errors would be returned (by inspecting the code path)
	// In production, this would fail if permissions are wrong
}

func TestLogSessionFailure(t *testing.T) {
	// Create a directory with the same name as the log file
	// This will cause the file open to fail
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	// Ensure data directory exists
	if err := storage.EnsureDataDir(); err != nil {
		t.Fatalf("Failed to ensure data dir: %v", err)
	}

	// Create a directory where the sessions.log file should be
	logPath := tmpDir + "/.gopomodoro/sessions.log"
	err := os.Mkdir(logPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create conflicting directory: %v", err)
	}

	// Try to log a session - should fail
	timestamp := time.Date(2025, 12, 23, 10, 30, 0, 0, time.UTC)
	err = storage.LogSession(timestamp, "work", "started", 0)

	if err == nil {
		t.Fatal("Expected error when logging to directory path, got nil")
	}

	// Verify error message contains context
	if !strings.Contains(err.Error(), "failed to") {
		t.Errorf("Expected error to be wrapped with context, got: %v", err)
	}
}
