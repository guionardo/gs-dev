package dev

import (
	"fmt"
	"strings"

	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/dialog"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/spf13/cobra"
)

type DevCommand struct {
	service *devservice.DevService
}

const (
	name        = "dev"
	description = "Rapid access to your development folders"
)

func (d *DevCommand) Setup(configuration *config.ConfigFile) error {
	d.service = devservice.NewService(configuration)
	return d.service.PurgeUnexistentRoots()
}

func (d *DevCommand) GetName() string {
	return name
}

func (d *DevCommand) GetDescription() string {
	return description
}

func (d *DevCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "arguments: <command> (can be add, delete, sync, list) or a query string to find a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				cmd.Help()
				return nil
			}

			return d.service.RunFind(args)
		},
		Annotations: map[string]string{
			consts.UseOutputAnnotation: "true",
		},
	}
	addCmd := &cobra.Command{
		Use:   "add",
		Short: "Add a folder to the roots",
		Long:  "arguments: <folder> (can be a relative, absolute and ~ resolved paths)",
		RunE:  d.addRoot,
		Args:  cobra.ExactArgs(1),
	}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a folder from the roots",
		Long:  "arguments: <folder> (can be a relative, absolute and ~ resolved paths)",
		RunE:  d.deleteRoot,
		Args:  cobra.ExactArgs(1),
	}
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: "Sync the roots",
		RunE:  d.syncRoots,
		Args:  cobra.ExactArgs(0),
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List the roots",
		RunE:  d.listRoots,
	}

	cmd.AddCommand(addCmd)
	cmd.AddCommand(deleteCmd)
	cmd.AddCommand(syncCmd)
	cmd.AddCommand(listCmd)

	return cmd
}

func (d *DevCommand) GetTUICommand() func() error {
	return func() error {
		command, err := dialog.Choose("dev", dialog.ToAnyArray([]string{"find", "add", "delete", "sync", "list"})...)
		if err != nil {
			return err
		}

		switch command {
		case "find":
			//TODO: Implement find command
			words, err := dialog.Input("Enter the words to find")
			if err != nil {
				return err
			}

			return d.service.RunFind(strings.Split(words, " "))

		case "add":
			root, err := dialog.InputDirectory("Enter the root folder", func(root string) error {
				root, err = d.service.CanAddRoot(root)
				if err != nil {
					return err
				}

				return nil
			})
			if err != nil {
				return err
			}

			if dialog.Confirm(colors.Parse("{primary}[%s]{normal} Are you sure you want to add this root?", root), true) {
				return d.service.AddRoot(root)
			}

			return nil
		case "delete":
			root, err := dialog.InputDirectory("Enter the root folder", func(root string) error {
				if !d.service.CanDeleteRoot(root) {
					return fmt.Errorf("root %s not registered", root)
				}

				return nil
			})
			if err != nil {
				return err
			}

			if dialog.Confirm(colors.Parse("{primary}[%s]{normal} Are you sure you want to delete this root?", root), true) {
				return d.service.DeleteRoot(root)
			}

			return nil
		case "sync":
			return d.service.Sync()
		case "list":
			return d.service.ListRoots()
		}

		fmt.Println(description)

		return nil
	}
}

func (d *DevCommand) addRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("Adding folder to roots")

	root, err := d.service.CanAddRoot(args[0])
	if err != nil {
		return err
	}

	return d.service.AddRoot(root)
}

func (d *DevCommand) deleteRoot(cmd *cobra.Command, args []string) error {
	fmt.Println("Deleting folder from roots")
	return nil
}

func (d *DevCommand) syncRoots(cmd *cobra.Command, args []string) error {
	return d.service.Sync()
}

func (d *DevCommand) listRoots(cmd *cobra.Command, args []string) error {
	return d.service.ListRoots()
}
