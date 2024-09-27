package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

func main() {

	taskManager := NewTaskManager("./tasks.json")
	taskManager.Init()

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
		Args:  cobra.MinimumNArgs(1),
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
		Args:  cobra.MinimumNArgs(1),
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
