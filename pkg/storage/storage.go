package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Can/sysrow/pkg/task"
)

// Storage handles task persistence
type Storage struct {
	mutex   sync.Mutex
	dataDir string
}

// TaskMetadata contains metadata about a task
type TaskMetadata struct {
	ID          string         `json:"id"`
	Command     string         `json:"command"`
	Status      task.TaskStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time     `json:"started_at,omitempty"`
	FinishedAt  *time.Time     `json:"finished_at,omitempty"`
	ExitCode    *int           `json:"exit_code,omitempty"`
	PID         *int           `json:"pid,omitempty"`
	GroupID     *string        `json:"group_id,omitempty"`
}

// NewStorage creates a new storage manager
func NewStorage(dataDir string) *Storage {
	return &Storage{
		dataDir: dataDir,
	}
}

// SaveTask saves a task's metadata to disk
func (s *Storage) SaveTask(t *task.Task) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create the tasks directory if it doesn't exist
	tasksDir := filepath.Join(s.dataDir, "tasks")
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Create task metadata
	metadata := TaskMetadata{
		ID:          t.ID,
		Command:     t.Command,
		Status:      t.Status,
		CreatedAt:   t.CreatedAt,
		ScheduledAt: t.ScheduledAt,
		StartedAt:   t.StartedAt,
		FinishedAt:  t.FinishedAt,
		ExitCode:    t.ExitCode,
		PID:         t.PID,
		GroupID:     t.GroupID,
	}

	// Marshal the metadata to JSON
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task metadata: %w", err)
	}

	// Write the metadata to file
	metadataPath := filepath.Join(tasksDir, t.ID+".json")
	if err := os.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return fmt.Errorf("failed to write task metadata file: %w", err)
	}

	return nil
}

// LoadTask loads a task's metadata from disk
func (s *Storage) LoadTask(id string) (*task.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Define the metadata file path
	metadataPath := filepath.Join(s.dataDir, "tasks", id+".json")

	// Read the metadata file
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task metadata file: %w", err)
	}

	// Unmarshal the metadata
	var metadata TaskMetadata
	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task metadata: %w", err)
	}

	// Create a task from the metadata
	t := &task.Task{
		ID:          metadata.ID,
		Command:     metadata.Command,
		Status:      metadata.Status,
		Priority:    task.PriorityNormal, // Default priority
		CreatedAt:   metadata.CreatedAt,
		ScheduledAt: metadata.ScheduledAt,
		StartedAt:   metadata.StartedAt,
		FinishedAt:  metadata.FinishedAt,
		ExitCode:    metadata.ExitCode,
		PID:         metadata.PID,
		GroupID:     metadata.GroupID,
	}

	return t, nil
}

// ListTasks returns all tasks
func (s *Storage) ListTasks() ([]*task.Task, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Define the tasks directory path
	tasksDir := filepath.Join(s.dataDir, "tasks")

	// Create the tasks directory if it doesn't exist
	if err := os.MkdirAll(tasksDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create tasks directory: %w", err)
	}

	// Read the tasks directory
	files, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	// Load each task
	tasks := make([]*task.Task, 0, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Extract the task ID from the filename
		taskID := file.Name()[:len(file.Name())-5] // Remove .json extension

		// Load the task
		task, err := s.LoadTask(taskID)
		if err != nil {
			// Log the error but continue loading other tasks
			fmt.Fprintf(os.Stderr, "Error loading task %s: %v\n", taskID, err)
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// DeleteTask deletes a task's metadata from disk
func (s *Storage) DeleteTask(id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Define the metadata file path
	metadataPath := filepath.Join(s.dataDir, "tasks", id+".json")

	// Delete the metadata file
	if err := os.Remove(metadataPath); err != nil {
		return fmt.Errorf("failed to delete task metadata file: %w", err)
	}

	return nil
}

// CleanupOldTasks removes tasks that are older than the specified duration
func (s *Storage) CleanupOldTasks(age time.Duration) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// List all tasks
	tasks, err := s.ListTasks()
	if err != nil {
		return fmt.Errorf("failed to list tasks: %w", err)
	}

	// Get the current time
	now := time.Now()

	// Delete old tasks
	for _, t := range tasks {
		// Skip tasks that are still running
		if t.Status == task.StatusRunning {
			continue
		}

		// Skip tasks that don't have a finish time
		if t.FinishedAt == nil {
			continue
		}

		// Check if the task is older than the specified age
		if now.Sub(*t.FinishedAt) > age {
			// Delete the task
			if err := s.DeleteTask(t.ID); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting task %s: %v\n", t.ID, err)
				continue
			}

			// Delete the task logs
			logsDir := filepath.Join(s.dataDir, "logs")
			logFiles := []string{
				filepath.Join(logsDir, t.ID+".stdout.log"),
				filepath.Join(logsDir, t.ID+".stderr.log"),
				filepath.Join(logsDir, t.ID+".app.log"),
			}

			for _, logFile := range logFiles {
				if _, err := os.Stat(logFile); err == nil {
					if err := os.Remove(logFile); err != nil {
						fmt.Fprintf(os.Stderr, "Error deleting log file %s: %v\n", logFile, err)
					}
				}
			}
		}
	}

	return nil
}
