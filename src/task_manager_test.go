package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

// Helper function to create a temporary TaskManager for testing
func createTempTaskManager(t *testing.T) *TaskManager {
	tempFile, err := os.CreateTemp("", "tasks_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFile.Close()

	tm := NewTaskManager(tempFile.Name())
	tm.Initialize()

	t.Cleanup(func() {
		os.Remove(tempFile.Name())
	})

	return tm
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestAddTask(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("Test Task 1")
	if tm.nextId != 2 {
		t.Errorf("Expected nextId to be 2, got %d", tm.nextId)
	}

	task, exists := tm.tasks[1]
	if !exists {
		t.Errorf("Task with ID 1 should exist")
	}

	if task.Description != "Test Task 1" {
		t.Errorf("Expected task description to be 'Test Task 1', got '%s'", task.Description)
	}

	if task.Completed {
		t.Errorf("New task should not be completed")
	}

	if !task.CompletedAt.IsZero() {
		t.Errorf("CompletedAt should be zero for new tasks")
	}

}

func TestDeleteTask(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("Task to delete 1") // ID 1
	tm.AddTask("Task to delete 2") // ID 2

	tm.DeleteTask(1)
	if _, exists := tm.tasks[1]; exists {
		t.Errorf("Task with ID 1 should have been deleted")
	}

	if len(tm.tasks) != 1 {
		t.Errorf("Expected 1 task after deletion, got %d", len(tm.tasks))
	}

	// attempt to delete non-existing task
	tm.DeleteTask(99)
	if len(tm.tasks) != 1 {
		t.Errorf("Deleting non-existing task should not change task count")
	}
}

func TestToggleTask(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("Task to toggle") // ID 1
	tm.ToggleTask(1)

	task := tm.tasks[1]

	if !task.Completed {
		t.Errorf("Task should be marked as completed")
	}

	if task.CompletedAt.IsZero() {
		t.Errorf("CompletedAt should be set when completed")
	}

	// toggle back to not completed
	tm.ToggleTask(1)

	task = tm.tasks[1]

	if task.Completed {
		t.Errorf("Task should be marked as not completed")
	}
	if !task.CompletedAt.IsZero() {
		t.Errorf("CompletedAt should be zero when not completed")
	}

	// attempt to toggle non-existing task
	tm.ToggleTask(99)
	if len(tm.tasks) != 1 {
		t.Errorf("Task count should remain 1 after attempting to toggle non-existing task")
	}
}

func TestClearTasks(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("Task 1") // ID 1 - not completed
	tm.AddTask("Task 2") // ID 2 - not completed
	tm.AddTask("Task 3") // ID 3 - not completed
	tm.ToggleTask(2)     // ID 2 - completed

	tm.ClearTasks(FilterCompleted)
	if len(tm.tasks) != 2 {
		t.Errorf("Expected 2 tasks after clearing completed tasks, got %d", len(tm.tasks))
	}
	if _, exists := tm.tasks[2]; exists {
		t.Errorf("Task with ID 2 should have been deleted")
	}

	tm.ClearTasks(FilterNotCompleted)
	if len(tm.tasks) != 0 {
		t.Errorf("Expected 0 tasks after clearing non-completed tasks, got %d", len(tm.tasks))
	}

	tm.AddTask("Task 4") // ID 4
	tm.AddTask("Task 5") // ID 5

	tm.ClearTasks(FilterAll)
	if len(tm.tasks) != 0 {
		t.Errorf("Expected 0 tasks after clearing all tasks, got %d", len(tm.tasks))
	}
}

func TestInitializeNextId(t *testing.T) {
	tempFile, err := os.CreateTemp("", "tasks_init_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}

	predefinedTasks := map[int]Task{
		1: {Description: "Predefined Task 1", Completed: false, CreatedAt: time.Now()},
		2: {Description: "Predefined Task 2", Completed: true, CreatedAt: time.Now(), CompletedAt: time.Now()},
		3: {Description: "Predefined Task 3", Completed: false, CreatedAt: time.Now()},
	}
	jsonData, err := json.MarshalIndent(predefinedTasks, "", " ")
	if err != nil {
		t.Fatalf("Failed to marshal predefined tasks: %v", err)
	}

	_, err = tempFile.Write(jsonData)
	if err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	tempFile.Close()

	tm := NewTaskManager(tempFile.Name())
	tm.Initialize()

	// nextId should be max existing ID + 1 = 4
	if tm.nextId != 4 {
		t.Errorf("Expected nextId to be 4, got %d", tm.nextId)
	}
}

func TestAddTaskEmptyDescription(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("")
	if len(tm.tasks) != 0 {
		t.Errorf("Adding a task with empty description should not create a task")
	}
}

func TestSaveTasksToFile(t *testing.T) {
	tm := createTempTaskManager(t)

	tm.AddTask("Task to save 1") // ID 1
	tm.AddTask("Task to save 2") // ID 2

	// read and verify contents
	fileData, err := os.ReadFile(tm.path)
	if err != nil {
		t.Fatalf("Failed to read tasks file: %v", err)
	}

	var loadedTasks map[int]Task
	err = json.Unmarshal(fileData, &loadedTasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal tasks from file: %v", err)
	}

	if len(loadedTasks) != 2 {
		t.Errorf("Expected 2 tasks in file, got %d", len(loadedTasks))
	}

	expectedDescriptions := map[int]string{
		1: "Task to save 1",
		2: "Task to save 2",
	}
	for id, desc := range expectedDescriptions {
		task, exists := loadedTasks[id]
		if !exists {
			t.Errorf("Task with ID %d should exist in file", id)
			continue
		}
		if task.Description != desc {
			t.Errorf("Task %d description mismatch. Expected '%s', got '%s'", id, desc, task.Description)
		}
	}
}
