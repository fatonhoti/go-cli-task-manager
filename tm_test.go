package main

import (
	"os"
	"testing"
)

func findTaskByID(tasks *map[int]Task, id int) (*Task, bool) {
	for _, task := range *tasks {
		if task.ID == id {
			return &task, true
		}
	}
	return nil, false
}

func TestAddTask(t *testing.T) {
	type testCase struct {
		testDesc    string // Description of the test case
		taskDesc    string // Description to add as a task
		expectExist bool   // Whether the task is expected to exist
		expectedID  int    // Expected ID of the task
	}

	// Define the test cases
	testCases := []testCase{
		{
			testDesc:    "Add task with valid description",
			taskDesc:    "Description for task 1",
			expectExist: true,
			expectedID:  1,
		},
		{
			testDesc:    "Add task with empty description",
			taskDesc:    "",
			expectExist: false, // empty descriptions are not allowed
			expectedID:  0,     // <-- zero-value for 'ID: int'
		},
	}

	// create temp file for testing to ensure isolation
	tempFile, err := os.CreateTemp("", "test_tasks_*.json")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	// make sure to remove the temporary file after the test
	defer os.Remove(tempFile.Name())

	// init TaskManager with the temporary file path
	tm := NewTaskManager(tempFile.Name())
	tm.Init()

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.testDesc, func(t *testing.T) {
			tm.AddTask(tc.taskDesc)

			task, exists := findTaskByID(&tm.tasks, tc.expectedID)

			if exists != tc.expectExist {
				t.Errorf("Task existence mismatch: expected %v, got %v", tc.expectExist, exists)
			}

			if exists {
				if task.ID != tc.expectedID {
					t.Errorf("Task ID mismatch: expected %d, got %d", tc.expectedID, task.ID)
				}
				if task.Description != tc.taskDesc {
					t.Errorf("Task Description mismatch: expected '%s', got '%s'", tc.taskDesc, task.Description)
				}
			}
		})
	}

	// verify that the tasks are saved correctly to the file
	tm2 := NewTaskManager(tempFile.Name())
	tm2.Init()

	if len(tm2.tasks) != len(tm.tasks) {
		t.Errorf("Loaded tasks count mismatch: expected %d, got %d", len(tm.tasks), len(tm2.tasks))
	}

	for _, task := range tm.tasks {
		loadedTask, exists := findTaskByID(&tm2.tasks, task.ID)
		if !exists {
			t.Errorf("Task with ID %d not found in loaded tasks", task.ID)
		} else {
			if loadedTask.Description != task.Description {
				t.Errorf("Loaded Task Description mismatch for ID %d: expected '%s', got '%s'", task.ID, task.Description, loadedTask.Description)
			}
			if loadedTask.Completed != task.Completed {
				t.Errorf("Loaded Task Completed status mismatch for ID %d: expected %v, got %v", task.ID, task.Completed, loadedTask.Completed)
			}
		}
	}
}
