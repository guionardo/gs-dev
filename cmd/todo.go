/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/guionardo/gs-dev/app/todo"
	"github.com/spf13/cobra"
)

// goCmd represents the go command
var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "To-do list",
	Long:  `Manage your tasks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("todo called")
	},
}

func init() {
	addCmd := &cobra.Command{
		Use:  "add",
		RunE: todo.RunAddTodo,
	}
	todoCmd.AddCommand(
		addCmd,
	)
	rootCmd.AddCommand(todoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// goCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// goCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
