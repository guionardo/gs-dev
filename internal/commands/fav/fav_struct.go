package fav

import (
	"context"
	"fmt"
	"io"

	contextdata "github.com/guionardo/gs-dev/internal/context"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
)

type (
	FavStruct struct {
		service *devservice.DevService
	}
)

func (fs *FavStruct) Run(ctx context.Context, output io.Writer) error {
	_, _ = fmt.Fprintln(output, "Running fav command")
	return fs.service.RunFavorites()
}

func (fs *FavStruct) Setup(ctx context.Context) error {
	cd, err := contextdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	fs.service = devservice.NewService(cd.RootConfig)

	return fs.service.PurgeUnexistentRoots()
}
