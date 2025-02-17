package build

import (
	"errors"
	"runtime/debug"

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
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return errors.New("unable to determine version information")
	}

	if buildInfo.Main.Version != "" {
		cmd.Printf("Version: %s\n", buildInfo.Main.Version)
	} else {
		cmd.Println("Version: unknown")
	}
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
