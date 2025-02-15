package main

import (
	"log/slog"
	"os"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/cmd"
	plugins_register "github.com/guionardo/gs-dev/plugins"
	"github.com/guionardo/gs-dev/plugins/manager"
)

// preRunArgs is a function that runs before the CLI command is executed,
func preRunArgs() []string {
	logLevel := slog.LevelInfo
	args := make([]string, 0, len(os.Args))

	for _, arg := range os.Args[1:] {
		if arg == "--debug" {
			logLevel = slog.LevelDebug
			continue
		}

		args = append(args, arg)
	}
	slog.SetLogLoggerLevel(logLevel)
	slog.Debug("Debug mode enabled")
	return args
}

func main() {
	args := preRunArgs()

	rootCmd := cmd.GetRootCmd(args, plugins_register.GetRegisteredPlugins(manager.Output)...)

	slog.Debug("Starting gs-dev", slog.String("version", build.Version))

	rootCmd.Execute()
}
