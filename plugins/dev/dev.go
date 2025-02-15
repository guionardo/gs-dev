package dev

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/guionardo/gs-dev/internal/arrays"
	"github.com/guionardo/gs-dev/internal/colors"
	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

const devName = "dev"

var (
	flagAdd  string
	flagDel  string
	flagSync bool
)

type DevPlugin struct {
	commons.BasePlugin
	configuration      Configuration
	rootsConfiguration RootsConfiguration
}

func Constructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &DevPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    devName,
			CanBeDisabled: true,
			Enabled:       true,
			Output:        output,
		},
	}
}

func (d *DevPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	d.BaseSetup(manager)

	if err = commons.LoadConfiguration(d, &d.configuration, false); err == nil {
		err = commons.LoadConfiguration(d, &d.rootsConfiguration, false, "roots")
	}
	if d.configuration.DefaultMaxDepth < 1 {
		d.configuration.DefaultMaxDepth = 3
		d.configuration.LastSync = time.Time{}
	}
	if int64(d.configuration.SyncInterval) == 0 {
		d.configuration.SyncInterval = time.Hour
		d.configuration.LastSync = time.Time{}
	}
	if d.rootsConfiguration.Roots == nil {
		d.rootsConfiguration.Roots = make(map[string]Root)
	}

	return err
}

func (d *DevPlugin) SetEnabled(enabled bool) error {
	d.configuration.Enabled = enabled
	d.Enabled = enabled
	return commons.SaveConfiguration(d, d.configuration)
}

func (d *DevPlugin) addRoot(root string) (err error) {
	if root, err = filepath.Abs(root); err != nil {
		return fmt.Errorf("error getting absolute path for %s: %w", root, err)
	}
	if stat, err := os.Stat(root); err != nil || !stat.IsDir() {
		return fmt.Errorf("error adding %s: not found or is not a directory", root)
	}
	if _, ok := d.rootsConfiguration.Roots[root]; ok {
		return fmt.Errorf("root %s already exists", root)
	}
	d.rootsConfiguration.Roots[root] = Root{MaxDepth: d.configuration.DefaultMaxDepth}
	return d.Sync()
}

func (d *DevPlugin) deleteRoot(root string) (err error) {
	if root, err = filepath.Abs(root); err != nil {
		return fmt.Errorf("error getting absolute path for %s: %w", root, err)
	}
	if _, ok := d.rootsConfiguration.Roots[root]; !ok {
		return fmt.Errorf("root %s not found", root)
	}
	delete(d.rootsConfiguration.Roots, root)
	return d.Sync()
}

func (d *DevPlugin) RunFind(cmd *cobra.Command, args []string) (err error) {
	slog.Debug("Running dev find", slog.Any("args", args))
	if flagAdd != "" {
		return d.addRoot(flagAdd)
	}
	if flagDel != "" {
		return d.deleteRoot(flagDel)
	}
	if flagSync {
		return d.Sync()
	}

	if d.shouldResync() {
		if err = d.Sync(); err != nil {
			return
		}
	}
	folders := d.Find(args)
	if len(folders) == 0 {
		return fmt.Errorf("no folders found for %v", args)
	}
	folder, err := chooseFolder(folders)

	if err == nil {

		d.WriteOutput("cd " + folder)

	}

	return err
}

func (d *DevPlugin) Find(words []string) []string {
	folders := make([]string, 0, 10)
	for rootFolder, root := range d.rootsConfiguration.Roots {
		for index := range root.Folders {
			if pathContainsPattern(root.Folders[index], rootFolder, words) {
				folders = append(folders, root.Folders[index])
			}
		}
	}
	return folders
}

func (d *DevPlugin) Sync() (err error) {
	filter := NewReaderFilter()
	for rootFolder, root := range d.rootsConfiguration.Roots {
		reader := NewDirReader(filter, rootFolder)
		folders := make([]string, 0, 16)
		for folder := range reader.Folders() {
			folders = append(folders, folder)
		}

		diffs := arrays.GetArraysDiffs(d.rootsConfiguration.Roots[rootFolder].Folders, folders)
		root.Folders = folders
		d.rootsConfiguration.Roots[rootFolder] = root
		if len(diffs) == 0 {
			colors.Normal("No changes\n")
		} else {
			colors.Blue("%d changes:\n", len(diffs))
			for index := range diffs {
				if diffs[index][0] == '+' {
					colors.Green(diffs[index] + "\n")
				} else {
					colors.Red(diffs[index] + "\n")
				}
			}
		}

	}
	d.configuration.LastSync = time.Now()
	if err = commons.SaveConfiguration(d, d.rootsConfiguration, "roots"); err == nil {
		err = commons.SaveConfiguration(d, d.configuration)
	}
	return
}

func (d *DevPlugin) RunSync(cmd *cobra.Command, args []string) error {
	slog.Debug("Running dev sync", slog.Any("args", args))
	return d.Sync()
}

func (d *DevPlugin) shouldResync() bool {
	return d.configuration.LastSync.Add(d.configuration.SyncInterval).Before(time.Now())
}

func (d *DevPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   devName,
		RunE:  d.RunFind,
		Short: "Rapid access to your development folders",
		Long: `This feature allows you to access your projects folder
quickly and get information about it`,
		Args: cobra.MinimumNArgs(1),
	}
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	cmd.Flags().StringVarP(&flagAdd, "add", "a", wd, "Add folder to roots")
	cmd.Flags().StringVarP(&flagDel, "delete", "d", wd, "Delete folder from roots")
	cmd.Flags().BoolVarP(&flagSync, "sync", "s", false, "Sync folders")

	return cmd
}
