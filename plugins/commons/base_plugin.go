package commons

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/generic"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
)

type BasePlugin struct {
	Manager       plugins.PluginsManager
	PluginName    string
	CanBeDisabled bool
	Enabled       bool
	Output        *outputfile.OutputFile
}

func (b BasePlugin) String() string {
	if !b.CanBeDisabled {
		return fmt.Sprintf("%s: allways enabled", b.PluginName)
	}
	return fmt.Sprintf("%s: %s", b.PluginName, generic.IfThen(b.Enabled, "enabled", "disabled"))
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
	if b.Output != nil {
		b.Output.AddContent(line)
	}
}
