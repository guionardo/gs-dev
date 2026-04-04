package gitstats

import (
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	gitstats "github.com/guionardo/gs-dev/internal/services/git"
	"github.com/spf13/cobra"
)

type GitStatsCommand struct {
	commands.CommonCommand

	service *gitstats.GitService
}

const (
	name        = "git-stats"
	description = "Show git stats"
)

func (g *GitStatsCommand) Init() {
	g.InitCommandVariables(name, description, false, "", false)
}
func (g *GitStatsCommand) Setup(configuration *config.ConfigFile) error {
	g.service = gitstats.NewGitService()
	return nil
}

func (g *GitStatsCommand) GetName() string {
	return name
}

func (g *GitStatsCommand) GetDescription() string {
	return description
}

func (g *GitStatsCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		Long:  "arguments: [<repositoryRoot>] (default: current directory)",
		RunE: func(cmd *cobra.Command, args []string) error {
			repositoryRoot := "."
			if len(args) > 0 {
				repositoryRoot = args[0]
			}

			return g.service.GetGitStats(repositoryRoot)
		},
	}
}

func (g *GitStatsCommand) GetTUICommand() func() error {
	return func() error {
		return g.service.GetGitStats(".")
	}
}
