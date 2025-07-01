package group

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/Can/sysrow/pkg/runner"
	"github.com/Can/sysrow/pkg/task"
	"github.com/google/uuid"
)

// Group represents a collection of related tasks
type Group struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	TaskIDs []string `json:"task_ids"`
}

// GroupManager manages task groups
type GroupManager struct {
	mutex     sync.Mutex
	dataDir   string
	groupsDir string
}

// NewGroupManager creates a new group manager
func NewGroupManager(dataDir string) *GroupManager {
	return &GroupManager{
		dataDir:   dataDir,
		groupsDir: filepath.Join(dataDir, "groups"),
	}
}

// CreateGroup creates a new task group
func (gm *GroupManager) CreateGroup(name string) (*Group, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Check if a group with the same name already exists
	groups, err := gm.ListGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	for _, g := range groups {
		if g.Name == name {
			return nil, fmt.Errorf("group with name %s already exists", name)
		}
	}

	// Create a new group
	group := &Group{
		ID:      uuid.New().String(),
		Name:    name,
		TaskIDs: make([]string, 0),
	}

	// Save the group
	if err := gm.saveGroup(group); err != nil {
		return nil, fmt.Errorf("failed to save group: %w", err)
	}

	return group, nil
}

// AddTask adds a task to a group
func (gm *GroupManager) AddTask(groupName string, command string) (*task.Task, error) {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Find the group
	group, err := gm.GetGroupByName(groupName)
	if err != nil {
		return nil, fmt.Errorf("failed to get group: %w", err)
	}

	// Create a new task
	t := task.NewTask(command, task.PriorityNormal)
	t.GroupID = &group.ID

	// Save the task
	if err := t.Save(); err != nil {
		return nil, fmt.Errorf("failed to save task: %w", err)
	}

	// Add the task ID to the group
	group.TaskIDs = append(group.TaskIDs, t.ID)

	// Save the updated group
	if err := gm.saveGroup(group); err != nil {
		return nil, fmt.Errorf("failed to save group: %w", err)
	}

	return t, nil
}

// RunGroup runs all tasks in a group
func (gm *GroupManager) RunGroup(groupName string, r *runner.Runner, background bool) error {
	// Find the group
	group, err := gm.GetGroupByName(groupName)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// Run each task in the group
	for _, taskID := range group.TaskIDs {
		// Load the task
		t, err := task.LoadTask(taskID)
		if err != nil {
			return fmt.Errorf("failed to load task %s: %w", taskID, err)
		}

		// Run the task
		if err := r.RunTask(t, background); err != nil {
			return fmt.Errorf("failed to run task %s: %w", taskID, err)
		}

		// If not running in background, wait for the task to complete
		if !background {
			fmt.Printf("Task %s completed with status: %s\n", t.ID, t.Status)
		}
	}

	return nil
}

// DeleteGroup deletes a group
func (gm *GroupManager) DeleteGroup(groupName string) error {
	gm.mutex.Lock()
	defer gm.mutex.Unlock()

	// Find the group
	group, err := gm.GetGroupByName(groupName)
	if err != nil {
		return fmt.Errorf("failed to get group: %w", err)
	}

	// Delete the group file
	groupPath := filepath.Join(gm.groupsDir, group.ID+".json")
	if err := os.Remove(groupPath); err != nil {
		return fmt.Errorf("failed to delete group file: %w", err)
	}

	return nil
}

// GetGroupByName returns a group by its name
func (gm *GroupManager) GetGroupByName(name string) (*Group, error) {
	// List all groups
	groups, err := gm.ListGroups()
	if err != nil {
		return nil, fmt.Errorf("failed to list groups: %w", err)
	}

	// Find the group with the matching name
	for _, g := range groups {
		if g.Name == name {
			return g, nil
		}
	}

	return nil, fmt.Errorf("group with name %s not found", name)
}

// ListGroups returns all groups
func (gm *GroupManager) ListGroups() ([]*Group, error) {
	// Create the groups directory if it doesn't exist
	if err := os.MkdirAll(gm.groupsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create groups directory: %w", err)
	}

	// Read the groups directory
	files, err := os.ReadDir(gm.groupsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read groups directory: %w", err)
	}

	// Load each group
	groups := make([]*Group, 0, len(files))
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// Extract the group ID from the filename
		groupID := file.Name()[:len(file.Name())-5] // Remove .json extension

		// Load the group
		group, err := gm.loadGroup(groupID)
		if err != nil {
			// Log the error but continue loading other groups
			fmt.Fprintf(os.Stderr, "Error loading group %s: %v\n", groupID, err)
			continue
		}

		groups = append(groups, group)
	}

	return groups, nil
}

// saveGroup saves a group to disk
func (gm *GroupManager) saveGroup(g *Group) error {
	// Create the groups directory if it doesn't exist
	if err := os.MkdirAll(gm.groupsDir, 0755); err != nil {
		return fmt.Errorf("failed to create groups directory: %w", err)
	}

	// Marshal the group to JSON
	groupData, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal group: %w", err)
	}

	// Write the group data to file
	groupPath := filepath.Join(gm.groupsDir, g.ID+".json")
	if err := os.WriteFile(groupPath, groupData, 0644); err != nil {
		return fmt.Errorf("failed to write group file: %w", err)
	}

	return nil
}

// loadGroup loads a group from disk by its ID
func (gm *GroupManager) loadGroup(id string) (*Group, error) {
	groupPath := filepath.Join(gm.groupsDir, id+".json")

	// Read the group file
	groupData, err := os.ReadFile(groupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read group file: %w", err)
	}

	// Unmarshal the group data
	var group Group
	if err := json.Unmarshal(groupData, &group); err != nil {
		return nil, fmt.Errorf("failed to unmarshal group: %w", err)
	}

	return &group, nil
}
