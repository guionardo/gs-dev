//go:build ignore

package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/cmd"
	"github.com/guionardo/gs-dev/internal/logging"
	plugins_register "github.com/guionardo/gs-dev/plugins"
	"github.com/spf13/cobra/doc"
)

func main() {
	var docsFolder string
	flag.StringVar(&docsFolder, "docs", "../../docs", "docs folder")
	flag.Parse()
	logging.Setup(false, os.Stdout)

	rootCmd := cmd.GetRootCmd([]string{}, plugins_register.GetRegisteredPlugins()...)
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
