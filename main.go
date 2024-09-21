package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

type Task struct {
	ID          int32     `json:"id"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"createdAt"`
	CompletedAt time.Time `json:"completedAt"`
}

type TaskManager struct {
	currMaxId int32
	tasks     []Task
	path      string
}

func (tm *TaskManager) Init() {
	tm.FromFile()
	var maxId int32 = 0
	for _, v := range tm.tasks {
		if v.ID > maxId {
			maxId = v.ID
		}
	}
	tm.currMaxId = maxId
}

func (tm *TaskManager) AddTask(desc string) {
	newTask := Task{
		tm.currMaxId + 1,
		desc,
		false,
		time.Now(),
		time.Time{},
	}
	tm.tasks = append(tm.tasks, newTask)
	tm.currMaxId++
}

func (tm *TaskManager) DeleteTask(id int32) {
	var idx int
	for i, task := range tm.tasks {
		if task.ID == id {
			idx = i
			break
		}
	}
	tm.tasks = append(tm.tasks[:idx], tm.tasks[idx+1:]...)
	// TODO: Shift IDs
}

func (tm *TaskManager) ListTasks() {
	for i := range tm.tasks {
		task := &tm.tasks[i]
		fmt.Printf("ID=%d\nDescription=%s\nCompleted=%t\nCreatedAt=%s\n", task.ID, task.Description, task.Completed, task.CreatedAt)
		if task.Completed {
			fmt.Printf("%s\n", task.CompletedAt)
		} else {
			fmt.Printf("CompletedAt=NOT_COMPLETED\n")
			fmt.Printf("----\n")
		}
	}
}

func (tm *TaskManager) CompleteTask(id int32) {
	for i := range tm.tasks {
		task := &tm.tasks[i]
		if task.ID == id {
			task.Completed = true
			task.CompletedAt = time.Now()
			break
		}
	}
}

func (tm *TaskManager) ToFile() {
	jsonData, err := json.Marshal(tm.tasks)
	if err != nil {
		fmt.Println("Failed. Could not serailize tasks.")
		panic(err)
	}

	f, err := os.OpenFile(tm.path, os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed. Could not open file %s\n", tm.path)
		panic(err)
	}

	_, err = f.Write(jsonData)
	if err != nil {
		fmt.Println("Failed. Could not write to file.")
		panic(err)
	}
}

func (tm *TaskManager) FromFile() {

	f, err := os.OpenFile(tm.path, os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}

	fi, err := f.Stat()
	if err != nil {
		panic(err)
	}

	// If we just created the file there will
	// be nothing to unmarshal.
	if fi.Size() > 0 {
		var buffer = make([]byte, fi.Size())
		f.Read(buffer)

		err = json.Unmarshal(buffer, &tm.tasks)
		if err != nil {
			panic(err)
		}
	}
}

func main() {

	var desc = flag.String("add", "", "Add a new task with a description.")
	var delete = flag.Int("delete", -271828, "Remove a specific task from the list.")
	var complete = flag.Int("complete", -271828, "Mark a specific task as completed.")
	var list = flag.Bool("list", false, "Display all tasks.")
	var help = flag.Bool("help", false, "Displays usage instructions.")

	flag.Parse()

	if *help {
		fmt.Println("Usage: go run main.go [flag]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	usedFlags := 0

	if *desc != "" {
		usedFlags++
	}

	if *delete != -271828 {
		usedFlags++
	}

	if *complete != -271828 {
		usedFlags++
	}

	if *list {
		usedFlags++
	}

	if usedFlags > 1 {
		fmt.Println("Error: Only one flag can be used at a time. Use '-help' to see usage instructions.")
		os.Exit(1)
	}

	taskManager := TaskManager{
		1,
		make([]Task, 0),
		"./tasks.json",
	}
	taskManager.Init()

	// Handle each flag case
	if *desc != "" {
		taskManager.AddTask(*desc)
	} else if *delete != -271828 {
		taskManager.DeleteTask(int32(*delete))
	} else if *complete != -271828 {
		taskManager.CompleteTask(int32(*complete))
	} else if *list {
		taskManager.ListTasks()
	} else {
		fmt.Println("No valid flag provided. Use -help for usage information.")
		os.Exit(1)
	}

	taskManager.ToFile()

	/*
		t1 := Task{
			1,
			"Some Description",
			false,
			time.Now(),
			time.Now(),
		}

		t2 := Task{
			2,
			"Some Other Description",
			true,
			time.Now(),
			time.Now(),
		}

		taskManager.AddTask(t1)
		taskManager.AddTask(t2)

		fmt.Println("Tasks:")
		taskManager.ListTasks()

		fmt.Println("Saving tasks to file...")
		taskManager.ToFile()

		fmt.Println("Loading tasks from file...")
		taskManager.FromFile()

		fmt.Println("Tasks after loading:")
		taskManager.ListTasks()
	*/

}
