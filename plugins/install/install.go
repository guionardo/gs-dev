package install

import (
	"fmt"
	"os"
	"strings"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/shell"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"

	"github.com/spf13/cobra"
)

type InstallShellPlugin struct {
	commons.BasePlugin

	sourceCommand string
}

func Constructor() plugins.CliPlugin {
	return &InstallShellPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    "install",
			CanBeDisabled: false,
			Enabled:       true,
		},
	}
}

func (i *InstallShellPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) error {
	i.BaseSetup(manager)

	// Source command
	//	source <(./go-dev init)
	executableName, err := os.Executable()
	if err != nil {
		return err
	}

	i.sourceCommand = fmt.Sprintf("source <(%s init)", executableName)

	return nil
}

func (i InstallShellPlugin) RunInstall(cmd *cobra.Command) error {
	//	source <(./gs-dev init)
	if strings.Contains(i.sourceCommand, "__debug_bin") {
		// Running from vscode
		return fmt.Errorf("bad executable name %s", i.sourceCommand)
	}

	profile, err := shell.NewProfileFile(build.AppName)
	if err != nil {
		return err
	}

	if line, ok := profile.HasEnabledCommandLine(); ok {
		return fmt.Errorf("binding was just installed into shell profile %s at line %d",
			profile.Path, line)
	}

	profile.SetFeature(i.sourceCommand, true)

	if err := profile.Save(); err != nil {
		return fmt.Errorf("error saving profile %s - %w", profile.Path, err)
	}

	cmd.Printf("binding installed at file %s\nbackup done at %s", profile.Path, profile.LastBackup)

	return nil
}

func (i InstallShellPlugin) RunUninstall(cmd *cobra.Command) error {
	profile, err := shell.NewProfileFile(build.AppName)
	if err != nil {
		return err
	}

	if _, ok := profile.HasEnabledCommandLine(); !ok {
		return fmt.Errorf("binding was not installed into shell profile %s",
			profile.Path)
	}

	profile.SetFeature(i.sourceCommand, false)

	if err := profile.Save(); err != nil {
		return fmt.Errorf("error saving profile %s - %w", profile.Path, err)
	}

	cmd.Printf("binding uninstalled at file %s\nbackup done at %s", profile.Path, profile.LastBackup)

	return nil
}

func (i InstallShellPlugin) Run(cmd *cobra.Command, args []string) error {
	if uninstall, err := cmd.Flags().GetBool("uninstall"); err != nil {
		return err
	} else {
		if i.sourceCommand == "" {
			return fmt.Errorf("cannot [un]install %s when running on vscode", build.AppName)
		}

		if uninstall {
			return i.RunUninstall(cmd)
		}

		return i.RunInstall(cmd)
	}
}
func (i InstallShellPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   i.PluginName,
		Short: "Install bindings on your shell profile",
		RunE:  i.Run,
	}
	cmd.Flags().BoolP("uninstall", "u", false, "Uninstall bindings")

	return cmd
}
