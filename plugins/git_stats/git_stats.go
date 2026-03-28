package gitstats

import (
	"fmt"

	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type GitStatsPlugin struct {
	commons.BasePlugin
}

func Constructor() plugins.CliPlugin {
	return &GitStatsPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName: "git-stats",
		},
	}
}

func (g *GitStatsPlugin) HasAlias() bool {
	return false
}

func (g *GitStatsPlugin) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "git-stats",
		Short: "Show git stats",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("git-stats")
		},
	}
}

func (g *GitStatsPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) error {
	return nil
}
