# SysRow

SysRow is a CLI-based background task manager tool designed for Linux environments. It allows users to queue, delay, group, and manage tasks with a modern and intuitive command-line interface.

## Features

- **Task Queueing**: Add tasks to a queue for sequential execution
- **Task Scheduling**: Delay tasks to run at specific times or after a duration
- **Task Grouping**: Create and manage groups of related tasks
- **Background Execution**: Run tasks in detached mode
- **Task Prioritization**: Assign priorities to tasks in the queue
- **Task Management**: List, check status, view logs, and cancel tasks

## Installation

```bash
# Clone the repository
git clone https://github.com/Can/sysrow.git

# Build the binary
cd sysrow
go build -o sysrow ./cmd/sysrow

# Install (optional)
mv sysrow /usr/local/bin/
```

## Usage

```bash
# Queue a task
sysrow queue "rsync -avz . backup:/data"

# Delay a task to run at a specific time
sysrow delay "npm run build" --at "02:00"

# Delay a task to run after a duration
sysrow delay "shutdown -r now" --after "5h"

# Create a task group
sysrow group create deploy

# Add tasks to a group
sysrow group add deploy "restart nginx"

# Run a task group
sysrow group run deploy

# Run a task in the background
sysrow run "heavy-process.sh" --bg

# Queue a task with priority
sysrow queue "render video.mp4" --priority high

# List all tasks
sysrow list

# Check task status
sysrow status <id>

# View task logs
sysrow logs <id>

# Cancel a task
sysrow cancel <id>

# Delete a task group
sysrow group delete deploy
```

## License

MIT
