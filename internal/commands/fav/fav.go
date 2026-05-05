package fav

import (
	"errors"

	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/spf13/cobra"
)

const (
	name        = "fav"
	description = "Rapid access to your last chosen folders"
)

type FavCommand struct {
	commands.CommonCommand

	service *devservice.DevService
}

func (f *FavCommand) Init() {
	f.InitCommandVariables(name, description, f, f.setup).WithInitAlias().WithOutput()
}

func (f *FavCommand) setup(configuration *config.ConfigRoot) error {
	f.service = devservice.NewService(configuration)
	return errors.Join(f.service.PurgeUnexistentRoots(), f.service.PurgeUnexistentFavorites())
}

func (f *FavCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&FavStruct{}, name, description, "", true)
}

func (f *FavCommand) GetTUICommand() func() error {
	return f.service.RunFavorites
}
