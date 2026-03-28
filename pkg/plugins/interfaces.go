package plugins

import (
	"iter"

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
		HasAlias() bool // if the plugin has an alias, it will be added to the init script
	}

	PluginsManager interface {
		Register(cp ...CliPlugin)
		Setup(configurationFolder string) error
		AddCommands(root *cobra.Command)

		GetPluginNames() []string
		GetPlugin(name string) (CliPlugin, bool)
		GetPlugins() iter.Seq[CliPlugin]

		GetRootCommand() *cobra.Command
	}

	PluginConstructor func() CliPlugin
)
