package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Logger handles logging for the application
type Logger struct {
	DataDir string
}

// NewLogger creates a new logger
func NewLogger(dataDir string) *Logger {
	return &Logger{
		DataDir: dataDir,
	}
}

// LogInfo logs an informational message
func (l *Logger) LogInfo(taskID, message string) error {
	return l.logMessage(taskID, "INFO", message)
}

// LogError logs an error message
func (l *Logger) LogError(taskID, message string) error {
	return l.logMessage(taskID, "ERROR", message)
}

// LogDebug logs a debug message
func (l *Logger) LogDebug(taskID, message string) error {
	return l.logMessage(taskID, "DEBUG", message)
}

// logMessage logs a message with the specified level
func (l *Logger) logMessage(taskID, level, message string) error {
	// Create the logs directory if it doesn't exist
	logsDir := filepath.Join(l.DataDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Define the log file path
	logPath := filepath.Join(logsDir, taskID+".app.log")

	// Open the log file in append mode
	logFile, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer logFile.Close()

	// Format the log message
	timestamp := time.Now().Format(time.RFC3339)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, level, message)

	// Write the log message to the file
	if _, err := logFile.WriteString(logLine); err != nil {
		return fmt.Errorf("failed to write to log file: %w", err)
	}

	return nil
}

// GetLogs returns the application logs for a task
func (l *Logger) GetLogs(taskID string) (string, error) {
	// Define the log file path
	logPath := filepath.Join(l.DataDir, "logs", taskID+".app.log")

	// Check if the log file exists
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		return "", nil
	}

	// Read the log file
	logData, err := os.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("failed to read log file: %w", err)
	}

	return string(logData), nil
}
