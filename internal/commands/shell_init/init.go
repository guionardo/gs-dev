package shellinit

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	installservice "github.com/guionardo/gs-dev/internal/services/install"
	"github.com/spf13/cobra"
)

type InitCommand struct {
	commands.CommonCommand

	installService *installservice.InstallService
}

const (
	name        = "init"
	description = "Initialization for shell alias"
)

func (i *InitCommand) Init() {
	i.InitCommandVariables(name, description, i, i.setup)
}

func (i *InitCommand) setup(configuration *config.ConfigRoot) error {
	i.installService = installservice.NewInstallService(configuration)

	return nil
}

func (i *InitCommand) BuildCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Print(i.installService.GetInitCommand())
			return nil
		},
	}
}

func (i *InitCommand) GetTUICommand() func() error {
	return func() error {
		fmt.Println("Run this command to initialize the shell alias")
		return nil
	}
}
