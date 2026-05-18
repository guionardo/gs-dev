//go:build ignore

package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/commands/dev"
	"github.com/guionardo/gs-dev/internal/commands/fav"
	gitstats "github.com/guionardo/gs-dev/internal/commands/git_stats"
	"github.com/guionardo/gs-dev/internal/commands/install"
	"github.com/guionardo/gs-dev/internal/commands/setup"
	shellinit "github.com/guionardo/gs-dev/internal/commands/shell_init"
	"github.com/guionardo/gs-dev/internal/commands/url"
	"github.com/guionardo/gs-dev/internal/logging"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	var docsFolder string
	flag.StringVar(&docsFolder, "docs", "../../docs", "docs folder")
	flag.Parse()
	logging.Setup(os.Stdout)

	rootCmd := getRootCommand()
	docsFolder, err := filepath.Abs(docsFolder)
	if err != nil {
		slog.Error("Failed to get docs folder", slog.String("folder", docsFolder), slog.Any("error", err))
		return
	}
	if err = os.RemoveAll(docsFolder); err != nil {
		slog.Error("Failed to remover docs folder", slog.String("folder", docsFolder), slog.Any("error", err))
		return
	}
	if err = pathtools.CreatePath(docsFolder); err != nil {
		slog.Error("Failed to create folder", slog.String("folder", docsFolder), slog.Any("error", err))
		return
	}
	if err = doc.GenMarkdownTree(rootCmd, docsFolder); err != nil {
		slog.Error("Failed to generate documentation", slog.String("folder", docsFolder), slog.Any("error", err))
		return
	}
	slog.Info("Documentation created", slog.String("folder", docsFolder))
}

func getRootCommand() *cobra.Command {
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
		&setup.SetupCommand{},
	)
	return commandsManager.GetRootCommand()
}
