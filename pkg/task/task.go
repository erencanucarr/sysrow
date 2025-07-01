package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// TaskStatus represents the current status of a task
type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
	StatusCancelled TaskStatus = "cancelled"
)

// TaskPriority represents the priority level of a task
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityNormal TaskPriority = "normal"
	PriorityHigh   TaskPriority = "high"
)

// Task represents a command to be executed
type Task struct {
	ID          string       `json:"id"`
	Command     string       `json:"command"`
	Status      TaskStatus   `json:"status"`
	Priority    TaskPriority `json:"priority"`
	CreatedAt   time.Time    `json:"created_at"`
	ScheduledAt *time.Time   `json:"scheduled_at,omitempty"`
	StartedAt   *time.Time   `json:"started_at,omitempty"`
	FinishedAt  *time.Time   `json:"finished_at,omitempty"`
	ExitCode    *int         `json:"exit_code,omitempty"`
	PID         *int         `json:"pid,omitempty"`
	GroupID     *string      `json:"group_id,omitempty"`
}

// DataDirectory is the path where all task data is stored
var DataDirectory string

// InitializeDataDirectory creates the necessary directory structure for storing task data
func InitializeDataDirectory() error {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	// Create the base data directory
	DataDirectory = filepath.Join(homeDir, ".sysrow")
	if err := os.MkdirAll(DataDirectory, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Create subdirectories
	dirs := []string{
		filepath.Join(DataDirectory, "tasks"),
		filepath.Join(DataDirectory, "groups"),
		filepath.Join(DataDirectory, "logs"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// NewTask creates a new task with the given command
func NewTask(command string, priority TaskPriority) *Task {
	taskID := uuid.New().String()

	return &Task{
		ID:        taskID,
		Command:   command,
		Status:    StatusPending,
		Priority:  priority,
		CreatedAt: time.Now(),
	}
}

// Save persists the task to disk
func (t *Task) Save() error {
	taskPath := filepath.Join(DataDirectory, "tasks", t.ID+".json")

	// Marshal the task to JSON
	taskData, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Write the task data to file
	if err := os.WriteFile(taskPath, taskData, 0644); err != nil {
		return fmt.Errorf("failed to write task file: %w", err)
	}

	return nil
}

// LoadTask loads a task from disk by its ID
func LoadTask(id string) (*Task, error) {
	taskPath := filepath.Join(DataDirectory, "tasks", id+".json")

	// Read the task file
	taskData, err := os.ReadFile(taskPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read task file: %w", err)
	}

	// Unmarshal the task data
	var task Task
	if err := json.Unmarshal(taskData, &task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// ListTasks returns all tasks
func ListTasks() ([]*Task, error) {
	tasksDir := filepath.Join(DataDirectory, "tasks")

	// Read the tasks directory
	files, err := os.ReadDir(tasksDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks directory: %w", err)
	}

	// Load each task
	tasks := make([]*Task, 0, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Extract the task ID from the filename
		taskID := file.Name()[:len(file.Name())-5] // Remove .json extension

		// Load the task
		task, err := LoadTask(taskID)
		if err != nil {
			// Log the error but continue loading other tasks
			fmt.Fprintf(os.Stderr, "Error loading task %s: %v\n", taskID, err)
			continue
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Cancel marks a task as cancelled
func (t *Task) Cancel() error {
	// Only pending or running tasks can be cancelled
	if t.Status != StatusPending && t.Status != StatusRunning {
		return fmt.Errorf("cannot cancel task with status %s", t.Status)
	}

	// If the task is running, we need to kill the process
	if t.Status == StatusRunning && t.PID != nil {
		// TODO: Implement process killing
	}

	// Update the task status
	t.Status = StatusCancelled
	now := time.Now()
	t.FinishedAt = &now

	// Save the updated task
	return t.Save()
}
