package install

import (
	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	installservice "github.com/guionardo/gs-dev/internal/services/install"
	"github.com/spf13/cobra"
)

type InstallCommand struct {
	service *installservice.InstallService
}

const (
	name        = "install"
	description = "Install the bindings into your shell profile"
)

func (i *InstallCommand) Setup(configuration *config.ConfigFile) error {
	i.service = installservice.NewInstallService(configuration)

	return nil
}

func (i *InstallCommand) GetName() string {
	return name
}

func (i *InstallCommand) GetDescription() string {
	return description
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
		if msg, isInstalled := i.service.IsInstalled(); isInstalled {
			colors.Primary(msg)
			if dialog.Confirm("Are you sure you want to uninstall the bindings?", false) {
				return i.service.Uninstall()
			}
		} else {
			colors.Secondary(msg)
			if dialog.Confirm("Are you sure you want to install the bindings?", true) {
				return i.service.Install()
			}
		}
		return nil
	}
}
