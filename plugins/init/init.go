package initshell

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/app/build"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

//go:embed init.sh
var initScript string

type InitShellPlugin struct {
	commons.BasePlugin
}

func Constructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &InitShellPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    "init",
			CanBeDisabled: false,
			Enabled:       true,
			Output:        output,
		},
	}
}

func (i InitShellPlugin) GetConfiguration() plugins.PluginConfiguration {
	return plugins.PluginConfiguration{
		Name:          "init",
		Enabled:       true,
		CanBeDisabled: false,
	}
}

func (i InitShellPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	i.BaseSetup(manager)
	return nil
}

func (i *InitShellPlugin) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   i.PluginName,
		Short: "Initialization for shell alias",
		Long: `Add to your profile script (.bashrc, etc)

		source <(gs-dev init)`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunInit()
		},
	}
}

func RunInit() error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}
	// #declare an array
	// my_array=("apple" "banana" "cherry" "date")

	doesntUseOutput := strings.Join([]string{`"url"`, `"init"`}, " ")
	output := path.Join(os.TempDir(), build.AppName)
	for key, value := range map[string]string{
		"GS_DEV":               executable,
		"GS_OUTPUT":            output,
		"GS_TOOL":              fmt.Sprintf("%s %s", build.AppName, build.Version),
		"GS_DOESNT_USE_OUTPUT": doesntUseOutput,
	} {
		initScript = strings.ReplaceAll(initScript, key, value)
	}

	fmt.Print(initScript)
	return nil
}
