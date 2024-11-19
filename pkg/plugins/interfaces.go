package plugins

import (
	"github.com/spf13/cobra"
)

type (
	PluginConfiguration struct {
		Name          string
		Enabled       bool
		CanBeDisabled bool
	}
	CliPlugin interface {
		Name() string
		Setup(manager PluginsManager, configurationFolder string) error
		GetConfiguration() PluginConfiguration

		IsEnabled() bool
		SetEnabled(enabled bool) error

		GetCobraCommand() *cobra.Command
	}

	PluginsManager interface {
		Register(cp ...CliPlugin)
		Setup(configurationFolder string) error
		AddCommands(root *cobra.Command)

		GetPluginNames() []string
		GetPlugin(name string) (CliPlugin, bool)

		GetRootCommand() *cobra.Command
	}
)
