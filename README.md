# 📝 Task Manager CLI

This is a simple task management CLI application built with Go, utilizing the [Cobra](https://github.com/spf13/cobra) package for command-line argument parsing and [tablewriter](https://github.com/olekukonko/tablewriter) for printing tables. The application allows users to **list**, **add**, **delete**, **toggle** (complete/pending), and **clear** tasks. The main goal of this project is to explore the basics of Go, with a secondary goal of developing a straightforward and easy to use task manager for personal use.

## Features

- **List tasks**: List all, completed, or non-completed tasks.
- **Add tasks**: Add new tasks to the task list.
- **Delete tasks**: Delete tasks by their ID.
- **Toggle tasks**: Toggle the completion status of tasks.
- **Clear tasks**: Clear all, completed, or non-completed tasks.

## Installation

### Prerequisites
- [Go](https://go.dev/dl/) version 1.22.5 (likely works with many other versions too)

### Building from Source
```bash
git clone https://github.com/fatonhoti/go-cli-task-manager.git
cd task-manager-cli
go build -o tm
```

## Usage demo

Once the application is built, you can use the `tm` executable to manage tasks. The following commands are available:

### 1. Add New Tasks
To add new tasks, run:

    ./tm add [task descriptions...]

Example:
    
    ./tm add "Buy groceries" "Finish homework" "Read a book"

Terminal Output:
```bash
Task added successfully. ID=1
Task added successfully. ID=2
Task added successfully. ID=3
```

### 2. List Tasks
To list tasks, use the `list` command with one of the following filters:
- `a`: List all tasks.
- `c`: List completed tasks.
- `nc`: List non-completed tasks.

Example:

    ./tm list a

Terminal Output:
```bash
+----+-----------+---------------------+---------------------+----------+-----------------+
| ID | Completed | Created At          | Completed At        | Days Ago | Description     |
+----+-----------+---------------------+---------------------+----------+-----------------+
| 1  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Buy groceries   |
| 2  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Finish homework |
| 3  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Read a book     |
+----+-----------+---------------------+---------------------+----------+-----------------+
```

Running:

    ./tm list c

Terminal Output:
    
    No tasks to show.

But now running

    ./tm list nc

yields:
```bash
+----+-----------+---------------------+---------------------+----------+-----------------+
| ID | Completed | Created At          | Completed At        | Days Ago | Description     |
+----+-----------+---------------------+---------------------+----------+-----------------+
| 1  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Buy groceries   |
| 2  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Finish homework |
| 3  | No        | 2024-09-27 14:25:43 | NOT_COMPLETED       | 0        | Read a book     |
+----+-----------+---------------------+---------------------+----------+-----------------+
```

### 3. Delete Tasks
To delete tasks by their ID, run:
    
    ./tm delete [task ids...]

Example:
    
    ./tm delete 1 3

Terminal Output:
```bash
Task 1 has been deleted.
Task 3 has been deleted.
```

### 4. Toggle Task Completion
To toggle the completion status of tasks, use:
    
    ./tm toggle [task ids...]

Example:
    
    ./tm toggle 2

Terminal Output:
    
    Task 2 has been marked completed.

### 5. Clear Tasks
To clear tasks based on their status, run:

    ./tm clear [a|c|nc]

- `a`: Clear all tasks.
- `c`: Clear completed tasks.
- `nc`: Clear non-completed tasks.

Example:
    
    ./tm clear c

Terminal Output:
    
    Cleared tasks.
