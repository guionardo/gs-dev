package dev

import (
	"fmt"
	"strings"

	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/dialog"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/spf13/cobra"
)

type DevCommand struct {
	commands.CommonCommand

	service *devservice.DevService
}

const (
	name        = "dev"
	description = "Rapid access to your development folders"

	addCommand    = "add"
	deleteCommand = "delete"
	syncCommand   = "sync"
	listCommand   = "list"
	findCommand   = "find"
	setupCommand  = "setup"

	addDescription    = "Add a folder to the roots"
	deleteDescription = "Delete a folder from the roots"
	syncDescription   = "Sync the roots"
	listDescription   = "List the roots"
	findDescription   = "Find a folder"
	setupDescription  = "Setup the current folder configuration"
)

func (d *DevCommand) Init() {
	d.InitCommandVariables(name, description, true, "", true)
}
func (d *DevCommand) Setup(configuration *config.ConfigFile) error {
	d.service = devservice.NewService(configuration)

	return d.service.PurgeUnexistentRoots()
}

func (d *DevCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "arguments: <command> (can be add, delete, sync, list) or a query string to find a folder",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}

			return d.service.RunFind(args)
		},
		Annotations: map[string]string{
			consts.UseOutputAnnotation: "true",
		},
	}
	addCmd := &cobra.Command{
		Use:   "add",
		Short: addDescription,
		Long:  "arguments: <folder> (can be a relative, absolute and ~ resolved paths)",
		RunE:  d.addRoot,
		Args:  cobra.ExactArgs(1),
	}
	deleteCmd := &cobra.Command{
		Use:   "delete",
		Short: deleteDescription,
		Long:  "arguments: <folder> (can be a relative, absolute and ~ resolved paths)",
		RunE:  d.deleteRoot,
		Args:  cobra.ExactArgs(1),
	}
	syncCmd := &cobra.Command{
		Use:   "sync",
		Short: syncDescription,
		RunE:  d.syncRoots,
		Args:  cobra.ExactArgs(0),
	}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: listDescription,
		RunE:  d.listRoots,
	}
	setupCmd := &cobra.Command{
		Use:   "setup",
		Short: setupDescription,
		RunE:  d.setup,
	}

	cmd.AddCommand(addCmd)
	cmd.AddCommand(deleteCmd)
	cmd.AddCommand(syncCmd)
	cmd.AddCommand(listCmd)
	cmd.AddCommand(setupCmd)

	return cmd
}

func (d *DevCommand) GetTUICommand() func() error {
	return func() error {
		command, err := dialog.Choose(
			"dev",
			dialog.ToAnyArray([]string{
				fmt.Sprintf("%s:%s", findCommand, findDescription),
				fmt.Sprintf("%s:%s", addCommand, addDescription),
				fmt.Sprintf("%s:%s", deleteCommand, deleteDescription),
				fmt.Sprintf("%s:%s", syncCommand, syncDescription),
				fmt.Sprintf("%s:%s", listCommand, listDescription),
			})...)
		if err != nil {
			return err
		}

		switch command {
		case findCommand:
			words, err := dialog.Input("Enter the words to find")
			if err != nil {
				return err
			}

			return d.service.RunFind(strings.Split(words, " "))

		case addCommand:
			root, err := dialog.InputDirectory("Enter the root folder", func(root string) error {
				_, err = d.service.CanAddRoot(root)
				if err != nil {
					return err
				}

				return nil
			})
			if err == nil && dialog.Confirm(colors.Parse("{primary}[%s]{normal} Are you sure you want to add this root?", root), true) {
				return d.service.AddRoot(root)
			}

			return err
		case deleteCommand:
			root, err := dialog.InputDirectory("Enter the root folder", func(root string) error {
				if !d.service.CanDeleteRoot(root) {
					return fmt.Errorf("root %s not registered", root)
				}

				return nil
			})
			if err == nil && dialog.Confirm(colors.Parse("{primary}[%s]{normal} Are you sure you want to delete this root?", root), true) {
				return d.service.DeleteRoot(root)
			}

			return err

		case syncCommand:
			return d.service.Sync()

		case listCommand:
			return d.service.ListRoots()
		}

		fmt.Println(description)

		return nil
	}
}

func (d *DevCommand) addRoot(cmd *cobra.Command, args []string) error {
	cmd.Println("Adding folder to roots")

	root, err := d.service.CanAddRoot(args[0])
	if err != nil {
		return err
	}

	return d.service.AddRoot(root)
}

func (d *DevCommand) deleteRoot(cmd *cobra.Command, args []string) error {
	cmd.Println("Deleting folder from roots")

	if !d.service.CanDeleteRoot(args[0]) {
		return fmt.Errorf("root %s not registered", args[0])
	}

	return d.service.DeleteRoot(args[0])
}

func (d *DevCommand) syncRoots(cmd *cobra.Command, args []string) error {
	return d.service.Sync()
}

func (d *DevCommand) listRoots(cmd *cobra.Command, args []string) error {
	return d.service.ListRoots()
}

func (d *DevCommand) setup(cmd *cobra.Command, args []string) error {
	return d.service.Setup()
}
