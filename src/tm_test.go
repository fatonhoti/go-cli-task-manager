package main

import (
	"os"
	"reflect"
	"testing"
)

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

			_, exists := tm.tasks[tc.expectedID]

			if exists != tc.expectExist {
				t.Errorf("Task existence mismatch: expected %v, got %v", tc.expectExist, exists)
			}
		})
	}

	// verify that the tasks are saved correctly to the file
	tm2 := NewTaskManager(tempFile.Name())
	tm2.Init()

	if len(tm2.tasks) != len(tm.tasks) {
		t.Errorf("Loaded tasks count mismatch: expected %d, got %d", len(tm.tasks), len(tm2.tasks))
	}

	if eq := reflect.DeepEqual(tm, tm2); !eq {
		t.Errorf("Loaded tasks differ from what was expected.")
	}

}
