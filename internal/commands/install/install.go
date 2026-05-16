package install

import (
	"fmt"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	installservice "github.com/guionardo/gs-dev/internal/services/install"
	"github.com/spf13/cobra"
)

type InstallCommand struct {
	commands.CommonCommand

	service *installservice.InstallService
}

const (
	name        = "install"
	description = "Install the bindings into your shell profile"
)

func (i *InstallCommand) Init() {
	i.InitCommandVariables(name, description, i, i.setup)
}
func (i *InstallCommand) setup(configuration *config.ConfigRoot) error {
	i.service = installservice.NewInstallService(configuration)

	return nil
}

func (i *InstallCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&InstallStruct{}, name, description, "", false)
}

func (i *InstallCommand) GetTUICommand() func() error {
	return func() error {
		if !build.IsValidBinary {
			return fmt.Errorf("cannot install/uninstall %s : %s", build.AppName, build.BuildInfo)
		}

		if msg, isInstalled := i.service.IsInstalled(); isInstalled {
			colors.Primary("%s", msg)

			if dialog.Confirm("Are you sure you want to uninstall the bindings?", false) {
				return i.service.Uninstall()
			}
		} else {
			colors.Secondary("%s", msg)

			if dialog.Confirm("Are you sure you want to install the bindings?", true) {
				return i.service.Install()
			}
		}

		return nil
	}
}
