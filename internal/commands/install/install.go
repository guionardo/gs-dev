package install

import (
	"fmt"

	"github.com/guionardo/gs-dev/app/build"
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
	i.InitCommandVariables(name, description, false, "", false)
}
func (i *InstallCommand) Setup(configuration *config.ConfigFile) error {
	i.service = installservice.NewInstallService(configuration)

	return nil
}

func (i *InstallCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			uninstall, err := cmd.Flags().GetBool("uninstall")
			if err != nil {
				return err
			}

			if uninstall {
				return i.service.Uninstall()
			}

			return i.service.Install()
		},
	}
	cmd.Flags().BoolP("uninstall", "u", false, "Uninstall bindings")

	return cmd
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
