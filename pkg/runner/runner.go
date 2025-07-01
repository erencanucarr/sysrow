package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/Can/sysrow/pkg/task"
)

// Runner is responsible for executing tasks
type Runner struct {
	DataDir string
}

// NewRunner creates a new task runner
func NewRunner(dataDir string) *Runner {
	return &Runner{
		DataDir: dataDir,
	}
}

// RunTask executes a task
func (r *Runner) RunTask(t *task.Task, background bool) error {
	// Update task status
	now := time.Now()
	t.Status = task.StatusRunning
	t.StartedAt = &now

	// Create log directory if it doesn't exist
	logsDir := filepath.Join(r.DataDir, "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %w", err)
	}

	// Create log files
	stdoutPath := filepath.Join(logsDir, t.ID+".stdout.log")
	stderrPath := filepath.Join(logsDir, t.ID+".stderr.log")

	stdoutFile, err := os.Create(stdoutPath)
	if err != nil {
		return fmt.Errorf("failed to create stdout log file: %w", err)
	}
	defer stdoutFile.Close()

	stderrFile, err := os.Create(stderrPath)
	if err != nil {
		return fmt.Errorf("failed to create stderr log file: %w", err)
	}
	defer stderrFile.Close()

	// Save the task state before execution
	if err := t.Save(); err != nil {
		return fmt.Errorf("failed to save task state: %w", err)
	}

	// Prepare the command
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", t.Command)
	} else {
		cmd = exec.Command("sh", "-c", t.Command)
	}

	// Set up command output
	cmd.Stdout = stdoutFile
	cmd.Stderr = stderrFile

	// Start the command
	if err := cmd.Start(); err != nil {
		// Update task status on error
		t.Status = task.StatusFailed
		endTime := time.Now()
		t.FinishedAt = &endTime
		exitCode := 1
		t.ExitCode = &exitCode
		t.Save()

		return fmt.Errorf("failed to start command: %w", err)
	}

	// Store the process ID
	pid := cmd.Process.Pid
	t.PID = &pid

	// Save the updated task state
	if err := t.Save(); err != nil {
		return fmt.Errorf("failed to save task state: %w", err)
	}

	// If running in background, return immediately
	if background {
		go func() {
			// Wait for the command to complete
			err := cmd.Wait()

			// Update task status
			endTime := time.Now()
			t.FinishedAt = &endTime
			t.PID = nil

			if err != nil {
				// Command failed
				t.Status = task.StatusFailed
				exitCode := 1
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				}
				t.ExitCode = &exitCode
			} else {
				// Command succeeded
				t.Status = task.StatusCompleted
				exitCode := 0
				t.ExitCode = &exitCode
			}

			// Save the final task state
			t.Save()
		}()

		return nil
	}

	// Wait for the command to complete
	err = cmd.Wait()

	// Update task status
	endTime := time.Now()
	t.FinishedAt = &endTime
	t.PID = nil

	if err != nil {
		// Command failed
		t.Status = task.StatusFailed
		exitCode := 1
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
		t.ExitCode = &exitCode
	} else {
		// Command succeeded
		t.Status = task.StatusCompleted
		exitCode := 0
		t.ExitCode = &exitCode
	}

	// Save the final task state
	if err := t.Save(); err != nil {
		return fmt.Errorf("failed to save task state: %w", err)
	}

	return nil
}

// GetTaskLogs returns the stdout and stderr logs for a task
func (r *Runner) GetTaskLogs(taskID string) (string, string, error) {
	// Define log file paths
	logsDir := filepath.Join(r.DataDir, "logs")
	stdoutPath := filepath.Join(logsDir, taskID+".stdout.log")
	stderrPath := filepath.Join(logsDir, taskID+".stderr.log")

	// Read stdout log
	stdoutData, err := os.ReadFile(stdoutPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read stdout log: %w", err)
	}

	// Read stderr log
	stderrData, err := os.ReadFile(stderrPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read stderr log: %w", err)
	}

	return string(stdoutData), string(stderrData), nil
}
