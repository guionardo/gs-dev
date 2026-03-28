package manager

import (
	"iter"
	"log/slog"
	"maps"
	"os"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/logging"
	"github.com/guionardo/gs-dev/pkg/plugins"
	postcommand "github.com/guionardo/gs-dev/pkg/post_command"
	"github.com/spf13/cobra"
)

type Manager struct {
	plugins map[string]plugins.CliPlugin
	rootCmd *cobra.Command
}

var Plugins plugins.PluginsManager

func init() {
	Plugins = &Manager{
		plugins: make(map[string]plugins.CliPlugin),
	}
}

func (p *Manager) Register(cp ...plugins.CliPlugin) {
	for _, plugin := range cp {
		p.plugins[plugin.Name()] = plugin
	}

	slog.Debug("Registered plugins", slog.Any("names", p.GetPluginNames()))
}

func (p *Manager) Setup(configurationFolder string) (err error) {
	for index := range p.plugins {
		err = p.plugins[index].Setup(p, configurationFolder)

		cfg := p.plugins[index].GetConfiguration()
		if err == nil {
			slog.Debug("Plugin setup", slog.String("plugin", cfg.Name), slog.Bool("enabled", p.plugins[index].IsEnabled()))
		} else {
			slog.Error("Plugin setup error", slog.String("plugin", cfg.Name), slog.Any("error", err))
			return
		}
	}

	return
}

func (p *Manager) AddCommands(root *cobra.Command) {
	for index := range p.plugins {
		if !p.plugins[index].IsEnabled() {
			continue
		}

		if cmd := p.plugins[index].GetCobraCommand(); cmd != nil {
			root.AddCommand(cmd)
		}
	}
}

func (p *Manager) GetPluginNames() []string {
	names := make([]string, len(p.plugins))

	index := 0
	for key := range p.plugins {
		names[index] = key
		index++
	}

	return names
}

func (p *Manager) GetPlugins() iter.Seq[plugins.CliPlugin] {
	return maps.Values(p.plugins)
}

func (p *Manager) GetPlugin(name string) (plugins.CliPlugin, bool) {
	if plugin, ok := p.plugins[name]; ok {
		return plugin, ok
	}

	return nil, false
}

func (p *Manager) PreRun(cmd *cobra.Command, args []string) {
	debug, _ := cmd.Flags().GetBool("debug")

	logging.PreSetupLog("Debug mode enabled", slog.LevelDebug)
	logging.Setup(debug, os.Stdout)
}

func (p *Manager) PostRun(cmd *cobra.Command, args []string) error {
	return postcommand.WriteOutput()
}

func (p *Manager) Run(cmd *cobra.Command, args []string) error {
	//TODO: default command is show the configuration of the plugins
	slog.Debug("Run command", slog.String("command", cmd.Name()), slog.Any("args", args))
	return nil
}

func (p *Manager) GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                build.AppName,
		Short:              build.ShortDescription,
		Long:               build.AppDescription,
		PersistentPreRun:   p.PreRun,
		PersistentPostRunE: p.PostRun,
		RunE:               p.Run,
	}
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug mode")

	for index := range p.plugins {
		if !p.plugins[index].IsEnabled() {
			continue
		}

		if cmd := p.plugins[index].GetCobraCommand(); cmd != nil {
			rootCmd.AddCommand(cmd)
		}
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%s %s\n", build.AppName, build.Version)
			cmd.Printf("Build info %s\n", build.BuildInfo)
		},
	}

	rootCmd.AddCommand(versionCmd)
	p.rootCmd = rootCmd

	return p.rootCmd
}
