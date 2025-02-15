package cmd

import (
	"log/slog"
	"os"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/manager"
	"github.com/spf13/cobra"
)

func setupPlugins() {
	manager.Plugins.Register()
}
func GetRootCmd(args []string, plugins ...plugins.CliPlugin) *cobra.Command {
	manager.Plugins.Register(plugins...)

	err := manager.Plugins.Setup(config.GetConfigDir())
	if err != nil {
		slog.Error("Error setting up plugins", slog.Any("error", err))
		os.Exit(1)
	}

	rootCmd := manager.Plugins.GetRootCommand()
	if cmd, _, err := rootCmd.Find(args); err != nil || cmd == nil {
		args = append([]string{"dev"}, args...)
	}
	rootCmd.SetArgs(args)
	return rootCmd
}
