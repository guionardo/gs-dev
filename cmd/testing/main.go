package main

import (
	"log"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/commands/dev"
	"github.com/guionardo/gs-dev/internal/commands/fav"
	gitstats "github.com/guionardo/gs-dev/internal/commands/git_stats"
	"github.com/guionardo/gs-dev/internal/commands/install"
	shellinit "github.com/guionardo/gs-dev/internal/commands/shell_init"
	"github.com/guionardo/gs-dev/internal/commands/url"
)

func main() {
	commandsManager, err := commands.NewCommandManager()
	if err != nil {
		log.Fatalf("Error creating commands manager: %v", err)
	}

	commandsManager.Register(
		&shellinit.InitCommand{},
		&dev.DevCommand{},
		&fav.FavCommand{},
		&gitstats.GitStatsCommand{},
		&install.InstallCommand{},
		&url.UrlCommand{},
	)
	rootCmd := commandsManager.GetRootCommand()

	_ = rootCmd.Execute()
}
