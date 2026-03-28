package commons

import (
	"fmt"

	"github.com/guionardo/go/flow"
	"github.com/guionardo/gs-dev/pkg/plugins"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
)

type BasePlugin struct {
	Manager       plugins.PluginsManager
	PluginName    string
	CanBeDisabled bool
	Enabled       bool
	AliasName     string
}

func (b BasePlugin) String() string {
	if !b.CanBeDisabled {
		return b.PluginName + ": allways enabled"
	}

	return fmt.Sprintf("%s: %s", b.PluginName, flow.If(b.Enabled, "enabled", "disabled"))
}

func (b BasePlugin) Name() string {
	return b.PluginName
}

func (b BasePlugin) IsEnabled() bool {
	return b.Enabled
}

func (b *BasePlugin) SetEnabled(enabled bool) error {
	if !enabled && !b.CanBeDisabled {
		return fmt.Errorf("plugin %s cannot be disabled", b.PluginName)
	}

	b.Enabled = enabled

	return nil
}

func (b *BasePlugin) BaseSetup(manager plugins.PluginsManager) {
	b.Manager = manager
}

func (b BasePlugin) GetConfiguration() plugins.PluginConfiguration {
	return plugins.PluginConfiguration{
		Name:          b.PluginName,
		Enabled:       b.Enabled,
		CanBeDisabled: b.CanBeDisabled,
	}
}

func (b BasePlugin) WriteOutput(line string) {
	postcommand.AddOutputLine(line)
}

func (b BasePlugin) HasAlias() bool {
	return b.AliasName != ""
}
