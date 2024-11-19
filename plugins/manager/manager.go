package manager

import (
	"log/slog"

	"github.com/guionardo/gs-dev/internal/metadata"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/spf13/cobra"
)

type Manager struct {
	plugins map[string]plugins.CliPlugin
	rootCmd *cobra.Command
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

func (p *Manager) GetPlugin(name string) (plugins.CliPlugin, bool) {
	if plugin, ok := p.plugins[name]; ok {
		return plugin.(plugins.CliPlugin), ok
	}
	return nil, false
}

func (p *Manager) GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{Use: metadata.AppName}
	rootCmd.Flags().Bool("debug", false, "Enable debug mode")
	for index := range p.plugins {
		if !p.plugins[index].IsEnabled() {
			continue
		}
		if cmd := p.plugins[index].GetCobraCommand(); cmd != nil {
			rootCmd.AddCommand(cmd)
		}
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%s %s\n", metadata.AppName, metadata.Version)
		},
	})
	p.rootCmd = rootCmd
	return p.rootCmd
}

var Plugins plugins.PluginsManager

func init() {
	Plugins = &Manager{
		plugins: make(map[string]plugins.CliPlugin),
	}
}
