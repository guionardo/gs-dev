package dev

import (
	"fmt"
	"strings"

	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	"github.com/guionardo/gs-dev/internal/interfaces"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/guionardo/gs-dev/pkg/console"
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

var _ interfaces.Command = &DevCommand{}

func (d *DevCommand) Init() {
	d.InitCommandVariables(name, description, d, d.setupFunc).WithDefaultArgument(findCommand).WithInitAlias().WithOutput()
}

func (d *DevCommand) setupFunc(configuration *config.ConfigRoot) error {
	d.service = devservice.NewService(configuration)

	return nil
}

func (d *DevCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&DevStruct{}, name, description, "", true)
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
			if err == nil && dialog.Confirm(console.Parse("{primary}[%s]{normal} Are you sure you want to add this root?", root), true) {
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
			if err == nil && dialog.Confirm(console.Parse("{primary}[%s]{normal} Are you sure you want to delete this root?", root), true) {
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
