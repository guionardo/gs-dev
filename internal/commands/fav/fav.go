package fav

import (
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/spf13/cobra"
)

const (
	name        = "fav"
	description = "Rapid access to your last chosen folders"
)

type FavCommand struct {
	service *devservice.DevService
}

func (f *FavCommand) Setup(configuration *config.ConfigFile) error {
	f.service = devservice.NewService(configuration)
	return f.service.PurgeUnexistentRoots()
}

func (f *FavCommand) GetName() string {
	return name
}

func (f *FavCommand) GetDescription() string {
	return description
}

func (f *FavCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			return f.service.RunFavorites()
		},
		Annotations: map[string]string{
			consts.UseOutputAnnotation: "true",
		},
	}
}

func (f *FavCommand) GetTUICommand() func() error {
	return func() error {
		return f.service.RunFavorites()
	}
}
