package main

import (
	"encoding/json"
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
	tasks []Task
	path  string
}

func (tm *TaskManager) AddTask(task Task) {
	tm.tasks = append(tm.tasks, task)
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
}

func (tm *TaskManager) ListTasks() {
	for _, task := range tm.tasks {
		fmt.Printf("ID=%d\nDescription=%s\nCompleted=%t\nCreatedAt=%s\nCompletedAt=<TODO>\n----\n", task.ID, task.Description, task.Completed, task.CompletedAt)
	}
}

func (tm *TaskManager) CompleteTask(id int32) {
	for _, task := range tm.tasks {
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

	taskManager := TaskManager{
		make([]Task, 0),
		"./tasks.json",
	}

	// Load tasks from disc
	fmt.Println("Loading tasks from disc...")
	taskManager.FromFile()

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

}
