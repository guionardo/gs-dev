package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/config"

	"github.com/guionardo/gs-dev/plugins/dev"
	initshell "github.com/guionardo/gs-dev/plugins/init"
	"github.com/guionardo/gs-dev/plugins/install"
	"github.com/guionardo/gs-dev/plugins/manager"
	"github.com/guionardo/gs-dev/plugins/plugins"
)

func setupPlugins() {
	manager.Plugins.Register(dev.NewDev(), plugins.NewPluginSetup(), initshell.NewInitShell(), install.NewInstallShell())
}

// preRunArgs is a function that runs before the CLI command is executed,
func preRunArgs() []string {
	logLevel := slog.LevelInfo
	args := make([]string, 0, len(os.Args))
	for i := range os.Args {
		if i == 0 {
			continue
		}
		arg := strings.ToLower(os.Args[i])
		if arg == "--debug" {
			logLevel = slog.LevelDebug
			continue
		}
		args = append(args, os.Args[i])
	}
	slog.SetLogLoggerLevel(logLevel)
	slog.Debug("Debug mode enabled")
	return args
}

func main() {
	args := preRunArgs()
	setupPlugins()
	slog.Debug("Starting gs-dev", slog.String("version", build.Version))

	err := manager.Plugins.Setup(config.GetConfigDir())
	if err != nil {
		slog.Error("Error setting up plugins", slog.Any("error", err))
		os.Exit(1)
	}

	rootCmd := manager.Plugins.GetRootCommand()
	rootCmd.SetArgs(args)
	if err = rootCmd.Execute(); err != nil {
		slog.Error("Error executing root command", slog.Any("error", err))
		os.Exit(1)
	}
}
