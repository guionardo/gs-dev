//go:build ignore

package main

import (
	"log/slog"
	"os"
	"path/filepath"

	"github.com/guionardo/gs-dev/internal/cmd"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	plugins_register "github.com/guionardo/gs-dev/plugins"
	"github.com/spf13/cobra/doc"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	output, err := outputfile.NewOutputFile("")
	rootCmd := cmd.GetRootCmd([]string{}, plugins_register.GetRegisteredPlugins(output)...)
	docsFolder, err := filepath.Abs("../docs")
	if err != nil {
		slog.Error("Failed to get ../docs folder", slog.Any("error", err))
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
	if err = doc.GenMarkdownTree(rootCmd, "../docs"); err != nil {
		slog.Error("Failed to generate documentation", slog.String("folder", docsFolder), slog.Any("error", err))
		return
	}
	slog.Info("Documentation created", slog.String("folder", docsFolder))
}
