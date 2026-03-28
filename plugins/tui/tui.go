package tui

import (
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type TUI struct {
	commons.BasePlugin
}

const tuiName = "tui"

func Constructor() plugins.CliPlugin {
	return &TUI{
		BasePlugin: commons.BasePlugin{
			PluginName:    tuiName,
			CanBeDisabled: true,
			Enabled:       true,
		},
	}
}

func (t *TUI) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   tuiName,
		RunE:  t.Run,
		Short: "Show TUI",
	}
}

func (t *TUI) Setup(manager plugins.PluginsManager, configurationFolder string) error {
	t.BaseSetup(manager)
	return nil
}

func (t *TUI) Run(cmd *cobra.Command, args []string) error {
	return nil
}
