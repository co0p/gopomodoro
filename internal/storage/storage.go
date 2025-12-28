package storage

import (
	"fmt"
	"os"
	"time"
)

const (
	dataDirName = ".gopomodoro"
	csvHeader   = "timestamp,session_type,event,duration_minutes"
)

// FileStorage implements the Storage interface for CSV file persistence
type FileStorage struct {
	logPath string
}

// NewFileStorage creates a new FileStorage instance
func NewFileStorage() *FileStorage {
	return &FileStorage{
		logPath: getSessionsLogPath(),
	}
}

// LogSession implements the Storage interface
func (fs *FileStorage) LogSession(timestamp time.Time, sessionType, status string, duration int) error {
	return logSessionToPath(fs.logPath, timestamp, sessionType, status, duration)
}

// EnsureDataDir creates the data directory if it doesn't exist
func EnsureDataDir() error {
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}
	return nil
}

// LogSession logs a session event to the CSV file (kept for backward compatibility)
func LogSession(timestamp time.Time, sessionType, event string, durationMinutes int) error {
	logPath := getSessionsLogPath()
	return logSessionToPath(logPath, timestamp, sessionType, event, durationMinutes)
}

// logSessionToPath writes a session log entry to the specified path
func logSessionToPath(logPath string, timestamp time.Time, sessionType, event string, durationMinutes int) error {
	// Check if file exists to determine if header is needed
	needsHeader := false
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		needsHeader = true
	}

	// Open file in append mode, create if doesn't exist
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open session log file: %w", err)
	}
	defer file.Close()

	// Write header if this is a new file
	if needsHeader {
		_, err = fmt.Fprintf(file, "%s\n", csvHeader)
		if err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
	}

	// Format and write CSV line
	csvLine := formatCSVLine(timestamp, sessionType, event, durationMinutes)
	_, err = fmt.Fprintf(file, "%s\n", csvLine)
	if err != nil {
		return fmt.Errorf("failed to write session log entry: %w", err)
	}
	return nil
}

// formatCSVLine formats a session entry as a CSV line
func formatCSVLine(timestamp time.Time, sessionType, event string, durationMinutes int) string {
	return fmt.Sprintf("%s,%s,%s,%d",
		timestamp.Format(time.RFC3339),
		sessionType,
		event,
		durationMinutes)
}

// getDataDir returns the full path to the data directory
func getDataDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home dir can't be determined
		return dataDirName
	}
	return homeDir + "/" + dataDirName
}

// getSessionsLogPath returns the full path to the sessions.log file
func getSessionsLogPath() string {
	return getDataDir() + "/sessions.log"
}
