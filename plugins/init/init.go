package initshell

import (
	_ "embed"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/guionardo/gs-dev/internal/metadata"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

//go:embed init.sh
var initScript string

type InitShell struct {
	commons.BasePlugin
}

func NewInitShell() *InitShell {
	return &InitShell{
		BasePlugin: commons.BasePlugin{
			PluginName:    "init",
			CanBeDisabled: false,
			Enabled:       true,
		},
	}
}

func (i InitShell) GetConfiguration() plugins.PluginConfiguration {
	return plugins.PluginConfiguration{
		Name:          "init",
		Enabled:       true,
		CanBeDisabled: false,
	}
}

func (i InitShell) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	i.BaseSetup(manager)
	return nil
}

func (i *InitShell) GetCobraCommand() *cobra.Command {
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
	output := path.Join(os.TempDir(), metadata.AppName)
	for key, value := range map[string]string{
		"GS_DEV":               executable,
		"GS_OUTPUT":            output,
		"GS_TOOL":              fmt.Sprintf("%s %s", metadata.AppName, metadata.Version),
		"GS_DOESNT_USE_OUTPUT": doesntUseOutput,
	} {
		initScript = strings.ReplaceAll(initScript, key, value)
	}

	fmt.Print(initScript)
	return nil
}
