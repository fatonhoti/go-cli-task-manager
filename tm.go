package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt"`
}

type TaskManager struct {
	nextId int
	tasks  map[int]Task
	path   string
}

func (tm *TaskManager) Init() {
	tm.LoadTasksFromFile()
	tm.nextId = 0
	for id := range tm.tasks {
		tm.nextId = max(tm.nextId, id)
	}
	tm.nextId++
}

func (tm *TaskManager) AddTask(desc string) {
	if len(desc) == 0 {
		return
	}

	tm.tasks[tm.nextId] = Task{
		Description: desc,
		Completed:   false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	tm.nextId++
	tm.SaveTasksToFile()
	fmt.Printf("Task added successfully. ID=%d\n", tm.nextId-1)
}

func (tm *TaskManager) DeleteTask(id int) {
	if _, ok := tm.tasks[id]; ok {
		delete(tm.tasks, id)
		tm.SaveTasksToFile()
		fmt.Printf("Task %d has been deleted.\n", id)
		return
	}
	fmt.Printf("error. No task with given ID=%d\n", id)
}

func (tm *TaskManager) CompleteTask(id int) {
	if task, ok := tm.tasks[id]; ok {
		task.Completed = true
		task.CompletedAt = time.Now()
		tm.tasks[id] = task
		tm.SaveTasksToFile()
		fmt.Printf("Task %d has been marked as completed.\n", id)
		return
	}
	fmt.Printf("error. No task with given ID=%d\n", id)
}

func (tm *TaskManager) UncheckTask(id int) {
	if task, ok := tm.tasks[id]; ok {
		task.Completed = false
		task.CompletedAt = time.Time{}
		tm.tasks[id] = task
		tm.SaveTasksToFile()
		fmt.Printf("Task %d has been marked as uncompleted.\n", id)
		return
	}
	fmt.Printf("error. No task with given ID=%d\n", id)
}

func (tm *TaskManager) ClearTasks(filter string) {
	ts := make(map[int]Task, 0)

	if filter == "completed" {
		for i := range tm.tasks {
			if !tm.tasks[i].Completed {
				ts[i] = tm.tasks[i]
			}
		}
	} else if filter == "pending" {
		for i := range tm.tasks {
			if tm.tasks[i].Completed {
				ts[i] = tm.tasks[i]
			}
		}
	}
	// if no filter is given => clear all tasks.

	tm.tasks = ts
	tm.SaveTasksToFile()
	fmt.Println("Cleared tasks.")
}

func (tm *TaskManager) ListTasks(viewCompact bool, filter string) {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks to show.")
		return
	}

	ts := make(map[int]Task, 0)

	if filter == "completed" {
		for id, task := range tm.tasks {
			if task.Completed {
				ts[id] = tm.tasks[id]
			}
		}
	} else if filter == "pending" {
		for id, task := range tm.tasks {
			if !task.Completed {
				ts[id] = tm.tasks[id]
			}
		}
	} else /* all */ {
		ts = tm.tasks
	}

	taskIds := make([]int, 0, len(ts))
	i := 0
	for k := range ts {
		taskIds = append(taskIds, k)
		i++
	}
	sort.Ints(taskIds)

	if viewCompact {
		tm.ListTasksCompact(&taskIds)
	} else {
		tm.ListTasksTable(&taskIds)
	}
}

func (tm *TaskManager) ListTasksCompact(taskIds *[]int) {
	fmt.Printf("--------------------------------------------\n")
	for id := range *taskIds {
		task := tm.tasks[id]
		status := "[ ]"
		if task.Completed {
			status = "[x]"
		}

		fmt.Printf("%s Task ID: %d\n", status, id)

		fmt.Println("Description:")
		for _, line := range strings.Split(task.Description, "\n") {
			fmt.Printf("%s\n", strings.Replace(line, `\n`, "\n", -1))
		}

		days := int(time.Since(task.CreatedAt).Hours() / 24)
		fmt.Printf("Created At: %s (%d days ago)\n", task.CreatedAt.Format("2006-01-02 15:04:05"), days)

		if task.Completed {
			fmt.Printf("Completed At: %s\n", task.CompletedAt.Format("2006-01-02 15:04:05"))
		}

		fmt.Printf("--------------------------------------------\n")
	}
}

func (tm *TaskManager) ListTasksTable(taskIds *[]int) {
	fmt.Printf("%-5s %-10s %-20s %-20s %-10s\n", "ID", "Completed", "Created At", "Completed At", "Days Ago")
	fmt.Println("-------------------------------------------------------------------")

	for _, id := range *taskIds {
		task := tm.tasks[id]
		completed := "No"
		completedAt := "NOT_COMPLETED"
		if task.Completed {
			completed = "Yes"
			completedAt = task.CompletedAt.Format("2006-01-02 15:04:05")
		}

		// Calculate how many days ago the task was created
		daysAgo := int(time.Since(task.CreatedAt).Hours() / 24)

		fmt.Printf("%-5d %-10s %-20s %-20s %-10d\n",
			id,
			completed,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			completedAt,
			daysAgo)

		// Print each line of the description
		for _, line := range strings.Split(task.Description, "\n") {
			fmt.Printf("%-s\n", strings.Replace(line, `\n`, "\n", -1))
		}

		fmt.Println("-------------------------------------------------------------------")
	}
}

func (tm *TaskManager) SaveTasksToFile() {

	// TODO: Consider adding recovery functionality in case below operations go wrong.

	f, err := os.OpenFile(tm.path, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error. Could not open file %s\n", tm.path)
		panic(err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(tm.tasks, "", " ")
	if err != nil {
		fmt.Println("error. Could not serailize tasks.")
		panic(err)
	}

	_, err = f.Write(jsonData)
	if err != nil {
		fmt.Println("error. Could not write to file.")
		panic(err)
	}
}

func (tm *TaskManager) LoadTasksFromFile() {

	f, err := os.OpenFile(tm.path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		fmt.Printf("Failed. Could not open file %s\n", tm.path)
		panic(err)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Printf("error. Stat() on %s failed.\n", tm.path)
		panic(err)
	}

	// If we just created the file there will
	// be nothing to unmarshal.
	if fi.Size() == 0 {
		return
	}

	var buffer = make([]byte, fi.Size())
	f.Read(buffer)

	err = json.Unmarshal(buffer, &tm.tasks)
	if err != nil {
		fmt.Println("error. Failed to unmarshal JSON.")
		panic(err)
	}

}

func NewTaskManager(path string) *TaskManager {
	return &TaskManager{
		nextId: 1,
		tasks:  map[int]Task{},
		path:   path,
	}
}

func main() {

	// list command
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	listCompact := listCmd.Bool("compact", false, "Display tasks in a compact format.")
	listFilter := listCmd.String("filter", "all", "List filtered tasks. Options are 'all', 'completed', or 'pending'.")

	// add command
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)

	// delete command
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	// complete command
	completeCmd := flag.NewFlagSet("complete", flag.ExitOnError)

	// uncheck (mark as uncompleted) command
	uncheckCmd := flag.NewFlagSet("uncheck", flag.ExitOnError)

	// clear command
	clearCmd := flag.NewFlagSet("clear", flag.ExitOnError)
	clearFilter := clearCmd.String("filter", "", "Clear tasks. Possible to filter, options are 'all', 'completed', or 'pending'.")

	if len(os.Args) < 2 {
		fmt.Println("error. Expected a command. Use 'help' to see usage instructions.")
		os.Exit(1)
	}

	// TODO: Consider allowing user to specify some other file via a CLI arg.
	taskManager := NewTaskManager("./tasks.json")
	taskManager.Init()

	switch os.Args[1] {
	case "list":
		listCmd.Parse(os.Args[2:])
		if *listFilter != "all" && *listFilter != "completed" && *listFilter != "pending" {
			fmt.Println("error. Invalid filter value.")
			listCmd.Usage()
			os.Exit(1)
		}
		taskManager.ListTasks(*listCompact, *listFilter)
	case "add":
		addCmd.Parse(os.Args[2:])
		args := addCmd.Args()
		if len(args) == 0 {
			fmt.Println("error. Must give at least one description corresponding to one task.")
			os.Exit(1)
		}
		for i := range args {
			taskManager.AddTask(args[i])
		}
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		args := deleteCmd.Args()
		if len(args) == 0 {
			fmt.Println("error. Must give at least one ID corresponding to one task.")
			os.Exit(1)
		}
		for i := range args {
			if id, err := strconv.Atoi(args[i]); err == nil {
				taskManager.DeleteTask(id)
			} else {
				fmt.Printf("error. Found non-integer argument ID=%s, skipping.", args[i])
			}
		}
	case "complete":
		completeCmd.Parse(os.Args[2:])
		args := completeCmd.Args()
		if len(args) == 0 {
			fmt.Println("error. Must give at least one ID corresponding to one task.")
			os.Exit(1)
		}
		for i := range args {
			if id, err := strconv.Atoi(args[i]); err == nil {
				taskManager.CompleteTask(id)
			} else {
				fmt.Printf("error. Found non-integer argument ID=%s, skipping.", args[i])
			}
		}
	case "uncheck":
		uncheckCmd.Parse((os.Args[2:]))
		args := uncheckCmd.Args()
		if len(args) == 0 {
			fmt.Println("error. Must give at least one ID corresponding to one task.")
			os.Exit(1)
		}
		for i := range args {
			if id, err := strconv.Atoi(args[i]); err == nil {
				taskManager.UncheckTask(id)
			} else {
				fmt.Printf("error. Found non-integer argument ID=%s, skipping.", args[i])
			}
		}
	case "clear":
		clearCmd.Parse(os.Args[2:])
		if *clearFilter != "all" && *clearFilter != "completed" && *clearFilter != "pending" {
			fmt.Println("error. Invalid filter value.")
			clearCmd.Usage()
			os.Exit(1)
		}
		taskManager.ClearTasks(*clearFilter)
	case "help":
		listCmd.Usage()
		fmt.Println("Usage of add:")
		fmt.Println("\t'tm add description1 description2 ...' creates a task per description.")
		fmt.Println("Usage of delete:")
		fmt.Println("\t'tm delete id1 id2 ...' deletes tasks with given IDs.")
		fmt.Println("Usage of complete:")
		fmt.Println("\t'tm complete id1 id2 ...' marks tasks with given IDs as completed.")
		fmt.Println("Usage of uncheck:")
		fmt.Println("\t'tm uncheck id1 id2 ...' marks tasks with given IDs as uncompleted.")
		clearCmd.Usage()
	default:
		fmt.Println("error. No subcommand found. Use 'help' to see usage instructions.")
		os.Exit(1)
	}

	os.Exit(0)
}
