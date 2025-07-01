package queue

import (
	"fmt"
	"sort"
	"sync"

	"github.com/Can/sysrow/pkg/task"
)

// Queue represents a task queue
type Queue struct {
	mutex sync.Mutex
	tasks []*task.Task
}

// NewQueue creates a new task queue
func NewQueue() *Queue {
	return &Queue{
		tasks: make([]*task.Task, 0),
	}
}

// Add adds a task to the queue
func (q *Queue) Add(t *task.Task) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Add the task to the queue
	q.tasks = append(q.tasks, t)

	// Sort the queue by priority
	q.sortByPriority()

	return nil
}

// sortByPriority sorts the queue by task priority
func (q *Queue) sortByPriority() {
	sort.Slice(q.tasks, func(i, j int) bool {
		// Define priority order: high > normal > low
		priorityOrder := map[task.TaskPriority]int{
			task.PriorityHigh:   3,
			task.PriorityNormal: 2,
			task.PriorityLow:    1,
		}

		return priorityOrder[q.tasks[i].Priority] > priorityOrder[q.tasks[j].Priority]
	})
}

// GetNext returns the next task in the queue
func (q *Queue) GetNext() *task.Task {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.tasks) == 0 {
		return nil
	}

	// Get the next task
	nextTask := q.tasks[0]

	// Remove the task from the queue
	q.tasks = q.tasks[1:]

	return nextTask
}

// Remove removes a task from the queue by ID
func (q *Queue) Remove(id string) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	for i, t := range q.tasks {
		if t.ID == id {
			// Remove the task from the queue
			q.tasks = append(q.tasks[:i], q.tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("task with ID %s not found in queue", id)
}

// List returns all tasks in the queue
func (q *Queue) List() []*task.Task {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	// Return a copy of the tasks slice to avoid race conditions
	tasksCopy := make([]*task.Task, len(q.tasks))
	copy(tasksCopy, q.tasks)

	return tasksCopy
}

// LoadQueue loads the queue from disk
func LoadQueue() (*Queue, error) {
	queue := NewQueue()

	// Load all pending tasks
	tasks, err := task.ListTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Add pending tasks to the queue
	for _, t := range tasks {
		if t.Status == task.StatusPending {
			queue.Add(t)
		}
	}

	return queue, nil
}

// SaveQueue saves the queue state to disk
func SaveQueue(q *Queue) error {
	// The queue state is implicitly saved through the tasks
	// Each task in the queue is already saved to disk
	return nil
}
