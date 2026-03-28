package gitstats

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/spf13/cobra"
)

type GitStatsCommand struct {
}

const (
	name        = "git-stats"
	description = "Show git stats"
)

func (g *GitStatsCommand) Setup(configuration *config.ConfigFile) error {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(description)
			return nil
		},
	}
}

func (g *GitStatsCommand) GetTUICommand() func() error {
	return func() error {
		fmt.Println(description)
		return nil
	}
}
