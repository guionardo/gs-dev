package shellinit

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/spf13/cobra"
)

type InitCommand struct {
}

const (
	name        = "init"
	description = "Initialization for shell alias"
)

func (i *InitCommand) Setup(configuration *config.ConfigFile) error {
	return nil
}

func (i *InitCommand) GetName() string {
	return name
}

func (i *InitCommand) GetDescription() string {
	return description
}

func (i *InitCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(description)
			return nil
		},
	}
}

func (i *InitCommand) GetTUICommand() func() error {
	return func() error {
		fmt.Println(description)
		return nil
	}
}
