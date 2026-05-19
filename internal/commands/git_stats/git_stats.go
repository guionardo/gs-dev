package gitstats

import (
	"github.com/guionardo/gs-dev/internal/cli"
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
	g.InitCommandVariables(name, description, g, g.setup)
}
func (g *GitStatsCommand) setup(configuration *config.ConfigRoot) error {
	g.service = gitstats.NewGitService()
	return nil
}

func (g *GitStatsCommand) GetName() string {
	return name
}

func (g *GitStatsCommand) GetDescription() string {
	return description
}
func (d *GitStatsCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&GitStatsStruct{}, name, description, "", false)
}

func (g *GitStatsCommand) GetTUICommand() func() error {
	return func() error {
		return g.service.GetGitStats(".", nil)
	}
}
