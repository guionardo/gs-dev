package manager

import (
	"log/slog"
	"os"
	"path"

	"github.com/guionardo/gs-dev/app/build"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/spf13/cobra"
)

type Manager struct {
	plugins map[string]plugins.CliPlugin
	rootCmd *cobra.Command
}

var flagOutput string

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
		return plugin, ok
	}
	return nil, false
}

func (p *Manager) PreRun(cmd *cobra.Command, args []string) error {
	Output.SetFile(flagOutput)
	return nil
}

func (p *Manager) PostRun(cmd *cobra.Command, args []string) error {
	Output.Close()
	return nil
}

func (p *Manager) GetRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:                build.AppName,
		Short:              build.ShortDescription,
		Long:               build.AppDescription,
		PersistentPreRunE:  p.PreRun,
		PersistentPostRunE: p.PostRun,
	}
	rootCmd.Flags().Bool("debug", false, "Enable debug mode")
	rootCmd.Flags().StringVarP(&flagOutput, "output", "o", path.Join(os.TempDir(), build.AppName), "Output script for shell alias")
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

var (
	Plugins plugins.PluginsManager
	Output  *outputfile.OutputFile
)

func init() {
	Plugins = &Manager{
		plugins: make(map[string]plugins.CliPlugin),
	}
	Output = &outputfile.OutputFile{}
}
