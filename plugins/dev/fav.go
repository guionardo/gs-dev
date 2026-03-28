package dev

import (
	"fmt"
	"log/slog"
	"os"
	"sort"

	"github.com/guionardo/gs-dev/internal/dialog"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type FavPlugin struct {
	commons.BasePlugin

	configuration Configuration
}

const favName = "fav"

func FavConstructor() plugins.CliPlugin {
	return &FavPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    favName,
			CanBeDisabled: true,
			Enabled:       true,
			AliasName:     "fav",
		},
	}
}

func (d *FavPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	d.BaseSetup(manager)

	if err = commons.LoadConfiguration(d, &d.configuration, false, "dev"); err != nil {
		return
	}

	if d.configuration.LastChosen == nil {
		d.configuration.LastChosen = make(LostChosenFolders)
	}

	return err
}

func (d *DevPlugin) ChosenFolder(folder string) {
	d.configuration.IncrementChosen(folder)
	_ = commons.SaveConfiguration(d, d.configuration)
}

func (d *FavPlugin) GetMostChosen(count int) (chosen []string) {
	type FC struct {
		folder string
		count  int
	}

	cc := make([]FC, 0, len(d.configuration.LastChosen))
	chosen = make([]string, 0, count)

	for folder, c := range d.configuration.LastChosen {
		if stat, err := os.Stat(folder); err == nil && stat.IsDir() {
			cc = append(cc, FC{
				folder,
				c,
			})
		}
	}

	sort.Slice(cc, func(i, j int) bool {
		return cc[i].count > cc[j].count
	})

	for i := 0; i < count && i < len(cc); i++ {
		chosen = append(chosen, cc[i].folder)
	}

	return
}

func (f *FavPlugin) RunChoose(cmd *cobra.Command, args []string) error {
	folders := f.GetMostChosen(f.configuration.MostChosenCount)
	slog.Debug("Getting most chosen", slog.Any("folders", folders))

	if len(folders) == 0 {
		return fmt.Errorf("no folders found for %v", args)
	}

	folder, err := dialog.Choose("Choose a folder:", (dialog.ToAnyArray(folders))...)
	slog.Debug("Chosen", slog.String("folder", folder))

	if err == nil {
		f.WriteOutput("cd " + folder)
	}

	return err
}

func (f *FavPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   favName,
		RunE:  f.RunChoose,
		Short: "Rapid access to your last chosen folders",
	}

	return cmd
}
