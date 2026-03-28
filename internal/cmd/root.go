package cmd

import (
	"log/slog"
	"os"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/pkg/plugins"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
	"github.com/guionardo/gs-dev/plugins/manager"
	"github.com/spf13/cobra"
)

func GetRootCmd(args []string, plugins ...plugins.CliPlugin) *cobra.Command {
	manager.Plugins.Register(plugins...)

	err := manager.Plugins.Setup(config.GetConfigDir())
	if err != nil {
		slog.Error("Error setting up plugins", slog.Any("error", err))
		os.Exit(1)
	}

	aliasCommands := make([]postcommand.Command, 0)

	for plugin := range manager.Plugins.GetPlugins() {
		if plugin.HasAlias() {
			aliasCommands = append(aliasCommands, postcommand.NewCommand(plugin.Name(), true, plugin.Name()))
		}
	}

	rootCmd := manager.Plugins.GetRootCommand()

	initSetup := postcommand.NewInitSetup(build.AppName, aliasCommands...)
	if err = initSetup.UpdateRootCommand(rootCmd); err != nil {
		slog.Error("Error updating root command", slog.Any("error", err))
		os.Exit(1)
	}

	// if no command is provided, use dev as default command
	if cmd, _, err := rootCmd.Find(args); err != nil || cmd == nil {
		args = append([]string{"dev"}, args...)
	}

	rootCmd.SetArgs(args)

	return rootCmd
}
