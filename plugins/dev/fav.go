package dev

import (
	"fmt"
	"log/slog"
	"os"
	"sort"

	outputfile "github.com/guionardo/gs-dev/internal/output_file"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

const favName = "fav"

type FavPlugin struct {
	commons.BasePlugin
	configuration Configuration
}

func FavConstructor(output *outputfile.OutputFile) plugins.CliPlugin {
	return &FavPlugin{
		BasePlugin: commons.BasePlugin{
			PluginName:    favName,
			CanBeDisabled: true,
			Enabled:       true,
			Output:        output,
		},
	}
}

func (d *FavPlugin) Setup(manager plugins.PluginsManager, configurationFolder string) (err error) {
	d.BaseSetup(manager)

	if err = commons.LoadConfiguration(d, &d.configuration, false, "dev"); err != nil {
		return
	}

	if d.configuration.LastChoosen == nil {
		d.configuration.LastChoosen = make(LastChoosenFolders)
	}

	return err
}

func (d *DevPlugin) ChoosedFolder(folder string) {
	d.configuration.IncrementChoosed(folder)
	commons.SaveConfiguration(d, d.configuration)
}

func (d *FavPlugin) GetMostChoosen(count int) (chosen []string) {
	type FC struct {
		folder string
		count  int
	}
	cc := make([]FC, 0, len(d.configuration.LastChoosen))
	chosen = make([]string, 0, count)
	for folder, c := range d.configuration.LastChoosen {
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
	folders := f.GetMostChoosen(5)
	slog.Debug("Getting most choosen", slog.Any("folders", folders))
	if len(folders) == 0 {
		return fmt.Errorf("no folders found for %v", args)
	}
	folder, err := chooseFolder(folders)
	slog.Debug("Choosed", slog.String("folder", folder))

	if err == nil {
		f.WriteOutput("cd " + folder)
	}

	return err

}

func (f *FavPlugin) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   favName,
		RunE:  f.RunChoose,
		Short: "Rapid access to your last choosen folders",
	}

	return cmd
}
