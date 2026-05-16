package main

import (
	"log"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/commands/dev"
	"github.com/guionardo/gs-dev/internal/commands/fav"
	gitstats "github.com/guionardo/gs-dev/internal/commands/git_stats"
	"github.com/guionardo/gs-dev/internal/commands/install"
	padcommand "github.com/guionardo/gs-dev/internal/commands/pad"
	"github.com/guionardo/gs-dev/internal/commands/setup"
	shellinit "github.com/guionardo/gs-dev/internal/commands/shell_init"
	"github.com/guionardo/gs-dev/internal/commands/url"
	"github.com/guionardo/gs-dev/internal/logging"
)

func main() {
	commandsManager, err := commands.NewCommandManager()
	if err != nil {
		log.Fatalf("Error creating commands manager: %v", err)
	}

	logging.Logger()

	commandsManager.Register(
		&shellinit.InitCommand{},
		&dev.DevCommand{},
		&fav.FavCommand{},
		&gitstats.GitStatsCommand{},
		&install.InstallCommand{},
		&url.UrlCommand{},
		&setup.SetupCommand{},
		&padcommand.PadCliCommand{},
		&padcommand.PadServerCommand{},
	)
	rootCmd := commandsManager.GetRootCommand()

	_ = rootCmd.Execute()
}
