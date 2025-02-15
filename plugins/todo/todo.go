package todo

import (
	"errors"

	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type TodoPlugin struct {
	commons.BasePlugin
}

const todoName = "todo"

func Constructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &TodoPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    todoName,
			CanBeDisabled: true,
			Enabled:       false,
			Output:        output,
		},
	}
}

func (t *TodoPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	//TODO: Implementar
	return nil
}
func (t *TodoPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   todoName,
		RunE:  t.Run,
		Short: "Use pad for exchange texts (bin.guiosoft.info)",
	}

	return cmd
}

func (t *TodoPlugin) Run(cmd *cobra.Command, args []string) error {
	return errors.New("unimplemented")
}
