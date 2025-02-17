package build

import (
	"github.com/guionardo/gs-dev/app/build"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type BuildInfoPlugin struct {
	commons.BasePlugin
}

const buildName = "build"

func Constructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &BuildInfoPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    buildName,
			CanBeDisabled: false,
			Enabled:       true,
			Output:        output,
		},
	}
}

func (b *BuildInfoPlugin) RunShowBuild(cmd *cobra.Command, args []string) error {
	cmd.Printf("Version: %s\n", build.Version)
	cmd.Printf("Build Info: %s\n", build.BuildInfo)
	return nil
}

func (b *BuildInfoPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   buildName,
		RunE:  b.RunShowBuild,
		Short: "Show build info",
	}

	return cmd
}

func (b *BuildInfoPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	return nil
}
