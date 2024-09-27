package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
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
}

func (tm *TaskManager) ToggleTask(id int) {
	if task, ok := tm.tasks[id]; ok {
		if task.Completed {
			task.Completed = false
			task.CompletedAt = time.Time{}
		} else {
			task.Completed = true
			task.CompletedAt = time.Now()
		}
		tm.tasks[id] = task
		tm.SaveTasksToFile()

		marked := "completed"
		if !task.Completed {
			marked = "pending"
		}
		fmt.Printf("Task %d has been marked %s.\n", id, marked)

		return
	}
}

func (tm *TaskManager) ClearTasks(filter string) {
	ts := make(map[int]Task, 0)

	if filter == "c" {
		for i := range tm.tasks {
			if tm.tasks[i].Completed {
				ts[i] = tm.tasks[i]
			}
		}
	} else if filter == "nc" {
		for i := range tm.tasks {
			if !tm.tasks[i].Completed {
				ts[i] = tm.tasks[i]
			}
		}
	}
	// if no filter is given => clear all tasks.

	tm.tasks = ts
	tm.SaveTasksToFile()
	fmt.Println("Cleared tasks.")
}

func (tm *TaskManager) ListTasks(filter string) {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks to show.")
		return
	}

	ts := make(map[int]Task, 0)

	if filter == "c" {
		for id, task := range tm.tasks {
			if task.Completed {
				ts[id] = tm.tasks[id]
			}
		}
	} else if filter == "nc" {
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

	fmt.Printf("%-5s %-10s %-20s %-20s %-10s\n", "ID", "Completed", "Created At", "Completed At", "Days Ago")
	fmt.Println("-------------------------------------------------------------------")

	for _, id := range taskIds {
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

	taskManager := NewTaskManager("./tasks.json")
	taskManager.Init()

	minId := 1
	maxId := 1
	for id := range taskManager.tasks {
		minId = min(minId, id)
		maxId = max(maxId, id)
	}

	// list command
	var cmdList = &cobra.Command{
		Use:   "list [a|c|nc]",
		Short: "List a set of tasks",
		Long:  "List a set of tasks. Use filter argument to filter for 'all', 'completed', or 'non-completed' tasks.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "a" && args[0] != "c" && args[0] != "nc" {
				fmt.Println("Error. Invalid filter.")
				return
			}
			taskManager.ListTasks(args[0])
		},
	}

	// add command
	var cmdAdd = &cobra.Command{
		Use:   "add [task descriptions...]",
		Short: "Add new tasks",
		Long:  "Add one or more new tasks with the specified descriptions to the task manager.",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			for _, description := range args {
				taskManager.AddTask(description)
			}
		},
	}

	// delete command
	var cmdDelete = &cobra.Command{
		Use:   "delete [task ids...]",
		Short: "Delete tasks by ID",
		Long:  "Delete one or more tasks by providing their respective task IDs. Only valid task IDs within the range are accepted.",
		Args:  cobra.RangeArgs(minId, maxId),
		Run: func(cmd *cobra.Command, args []string) {
			for _, id := range args {
				if n, err := strconv.Atoi(id); err == nil {
					taskManager.DeleteTask(n)
				}
			}
		},
	}

	// toggle command
	var cmdToggle = &cobra.Command{
		Use:   "toggle [task ids...]",
		Short: "Toggle the completion status of tasks",
		Long:  "Toggle the completion status of one or more tasks by their task IDs. This changes completed tasks to non-completed and vice versa.",
		Args:  cobra.RangeArgs(minId, maxId),
		Run: func(cmd *cobra.Command, args []string) {
			for _, id := range args {
				if n, err := strconv.Atoi(id); err == nil {
					taskManager.ToggleTask(n)
				}
			}
		},
	}

	// clear command
	var cmdClear = &cobra.Command{
		Use:   "clear [a|c|nc]",
		Short: "Clear tasks based on completion status",
		Long:  "Clear tasks by providing a filter: 'a' for all tasks, 'c' for completed tasks, or 'nc' for non-completed tasks.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if args[0] != "a" && args[0] != "c" && args[0] != "nc" {
				fmt.Println("Error. Invalid filter.")
				return
			}
			taskManager.ClearTasks(args[0])
		},
	}

	var rootCmd = &cobra.Command{Use: "tm"}

	rootCmd.AddCommand(cmdList)
	rootCmd.AddCommand(cmdAdd)
	rootCmd.AddCommand(cmdDelete)
	rootCmd.AddCommand(cmdToggle)
	rootCmd.AddCommand(cmdClear)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

	os.Exit(0)
}
