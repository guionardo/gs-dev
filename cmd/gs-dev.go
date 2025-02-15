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
	"github.com/guionardo/gs-dev/plugins/url"
)

func setupPlugins() {
	manager.Plugins.Register(dev.NewDev(), plugins.NewPluginSetup(), initshell.NewInitShell(), install.NewInstallShell(), url.NewURL())
}

// preRunArgs is a function that runs before the CLI command is executed,
func preRunArgs() ([]string, string) {
	logLevel := slog.LevelInfo
	args := make([]string, 0, len(os.Args))
	var command string
	for _, arg := range os.Args[1:] {
		if arg == "--debug" {
			logLevel = slog.LevelDebug
			continue
		}
		if len(command) == 0 && !strings.HasPrefix(arg, "--") {
			command = arg
		}
		args = append(args, arg)
	}
	slog.SetLogLoggerLevel(logLevel)
	slog.Debug("Debug mode enabled")
	return args, command
}

func main() {
	args, _ := preRunArgs()
	setupPlugins()
	slog.Debug("Starting gs-dev", slog.String("version", build.Version))

	err := manager.Plugins.Setup(config.GetConfigDir())
	if err != nil {
		slog.Error("Error setting up plugins", slog.Any("error", err))
		os.Exit(1)
	}

	rootCmd := manager.Plugins.GetRootCommand()
	if cmd, _, err := rootCmd.Find(args); err != nil || cmd == nil {
		args = append([]string{"dev"}, args...)
	}
	rootCmd.SetArgs(args)

	rootCmd.Execute()
}
