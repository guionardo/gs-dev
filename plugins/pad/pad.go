package pad

import (
	"errors"

	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type PadPlugin struct {
	commons.BasePlugin
}

const padName = "pad"

func Constructor() plugins.CliPlugin {
	return &PadPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    padName,
			CanBeDisabled: true,
			Enabled:       false,
		},
	}
}

func (p *PadPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	//TODO: Implementar
	return nil
}

func (p *PadPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   padName,
		RunE:  p.Run,
		Short: "Use pad for exchange texts (bin.guiosoft.info)",
	}

	return cmd
}

func (p *PadPlugin) Run(cmd *cobra.Command, args []string) error {
	return errors.New("unimplemented")
}
