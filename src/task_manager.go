package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
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
	if _, ok := tm.tasks[id]; !ok {
		fmt.Printf("Task %d not found.\n", id)
		return
	}

	delete(tm.tasks, id)
	tm.SaveTasksToFile()
	fmt.Printf("Task %d has been deleted.\n", id)
}

func (tm *TaskManager) ToggleTask(id int) {
	task, ok := tm.tasks[id]

	if !ok {
		fmt.Printf("Task %d not found.\n", id)
		return
	}

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
}

func (tm *TaskManager) ClearTasks(filter string) {
	ts := make(map[int]Task, 0)

	for id, task := range tm.tasks {
		if (filter == "c" && !task.Completed) || (filter == "nc" && task.Completed) {
			ts[id] = task
		}
	}

	if filter != "c" && filter != "nc" {
		ts = make(map[int]Task) // clear all
	}

	tm.tasks = ts
	tm.SaveTasksToFile()

	fmt.Println("Cleared tasks.")
}

func (tm *TaskManager) ListTasks(filter string) {
	if len(tm.tasks) == 0 {
		fmt.Println("No tasks to show.")
		return
	}

	ts := make(map[int]Task)
	for id, task := range tm.tasks {
		if (filter == "c" && task.Completed) || (filter == "nc" && !task.Completed) {
			ts[id] = task
		}
	}

	if filter != "c" && filter != "nc" {
		ts = tm.tasks // all tasks.
	}

	// sort by task ID because Go maps are not sorted by default
	taskIds := make([]int, 0, len(ts))
	for id := range ts {
		taskIds = append(taskIds, id)
	}
	sort.Ints(taskIds)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Completed", "Created At", "Completed At", "Days Ago", "Description"})

	for _, id := range taskIds {
		task := ts[id]

		completed := "No"
		completedAt := "NOT_COMPLETED"
		if task.Completed {
			completed = "Yes"
			completedAt = task.CompletedAt.Format("2006-01-02 15:04:05")
		}

		daysAgo := int(time.Since(task.CreatedAt).Hours() / 24)

		table.Append([]string{
			fmt.Sprintf("%d", id),
			completed,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			completedAt,
			fmt.Sprintf("%d", daysAgo),
			task.Description,
		})
	}

	table.Render()
}

func (tm *TaskManager) SaveTasksToFile() {

	tmpFilePath := tm.path + ".tmp"
	f, err := os.OpenFile(tmpFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Printf("error. Could not open temporary file %s\n", tmpFilePath)
		panic(err)
	}
	defer f.Close()

	jsonData, err := json.MarshalIndent(tm.tasks, "", " ")
	if err != nil {
		fmt.Println("error. Could not serialize tasks.")
		panic(err)
	}

	_, err = f.Write(jsonData)
	if err != nil {
		fmt.Println("error. Could not write to temporary file.")
		panic(err)
	}

	if err := f.Close(); err != nil {
		fmt.Println("error. Could not close temporary file.")
		panic(err)
	}

	err = os.Rename(tmpFilePath, tm.path)
	if err != nil {
		fmt.Println("error. Could not replace the original file with the temporary file.")
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
