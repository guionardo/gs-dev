package plugins_setup

import (
	"fmt"
	"strings"

	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type PluginSetup struct {
	commons.BasePlugin
}

func Constructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &PluginSetup{
		BasePlugin: commons.BasePlugin{
			PluginName:    "plugins",
			CanBeDisabled: false,
			Enabled:       true,
			Output:        output,
		},
	}

}

func (d *PluginSetup) GetConfiguration() plugins.PluginConfiguration {
	return plugins.PluginConfiguration{
		Name:          d.PluginName,
		Enabled:       true,
		CanBeDisabled: false,
	}
}

func (d *PluginSetup) IsEnabled() bool {
	return true
}
func (d *PluginSetup) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	d.BaseSetup(manager)
	return nil
}

func (d *PluginSetup) SetEnabled(enabled bool) error {
	if !enabled {
		return fmt.Errorf("plugin %s cannot be disabled", d.PluginName)
	}
	return nil
}

func (d *PluginSetup) RunList(cmd *cobra.Command, args []string) error {
	cmd.Printf("Plugins\n")
	for _, name := range d.Manager.GetPluginNames() {
		plugin, _ := d.Manager.GetPlugin(name)
		cfg := plugin.GetConfiguration()
		cmd.Printf(" %s = %v\n", cfg.Name, cfg.Enabled)
	}
	return nil
}

func (d *PluginSetup) runEnable(cmd *cobra.Command, pluginName string, enabled bool) (err error) {
	if plugin, ok := d.Manager.GetPlugin(pluginName); ok {
		if err = plugin.SetEnabled(enabled); err != nil {
			return
		}
		if err = plugin.SetEnabled(enabled); err == nil {
			cmd.Printf("Plugin %s enabled: %v\n", pluginName, enabled)
		}
	} else {
		err = fmt.Errorf("plugin %s not found. use one of these [%s]", pluginName, strings.Join(d.Manager.GetPluginNames(), ", "))
	}
	if err != nil {
		cmd.Printf("Error setting enable = %v on plugin %s: %v\n", enabled, pluginName, err)
	}
	return
}
func (d *PluginSetup) RunEnable(cmd *cobra.Command, args []string) error {
	return d.runEnable(cmd, args[0], true)
}
func (d *PluginSetup) RunDisable(cmd *cobra.Command, args []string) error {
	return d.runEnable(cmd, args[0], false)
}

func (d *PluginSetup) GetCobraCommand() *cobra.Command {
	pluginsNames := strings.Join(d.Manager.GetPluginNames(), " | ")
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List plugins",
		RunE:  d.RunList,
	}

	enableCmd := &cobra.Command{
		Use:   fmt.Sprintf("enable [%s]", pluginsNames),
		Short: "Enable plugin",
		Args:  cobra.ExactArgs(1),
		RunE:  d.RunEnable,
	}

	disableCmd := &cobra.Command{
		Use:   fmt.Sprintf("disable [%s]", pluginsNames),
		Short: "Disable plugin",
		Args:  cobra.ExactArgs(1),
		RunE:  d.RunDisable,
	}

	cmd := &cobra.Command{
		Use:   d.PluginName,
		Short: "Setup plugins",
	}
	cmd.AddCommand(listCmd, enableCmd, disableCmd)
	return cmd
}
