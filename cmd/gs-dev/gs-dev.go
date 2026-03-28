package main

import (
	"log/slog"
	"os"

	_ "github.com/spf13/cobra/doc"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/cmd"
	plugins_register "github.com/guionardo/gs-dev/plugins"
)

func main() {
	rootCmd := cmd.GetRootCmd(os.Args[1:], plugins_register.GetRegisteredPlugins()...)

	slog.Debug("Starting gs-dev", slog.String("version", build.Version))

	_ = rootCmd.Execute()
}
