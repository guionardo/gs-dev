package dev

import (
	"context"

	"fmt"
	"io"

	contextdata "github.com/guionardo/gs-dev/internal/context"
	"github.com/guionardo/gs-dev/internal/errors"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
)

type (
	baseDevStruct struct {
		service *devservice.DevService
	}

	DevStruct struct {
		baseDevStruct

		Add      AddCommandStruct    `subcommand:"add,add a folder to the roots"`
		Delete   DeleteCommandStruct `subcommand:"delete,delete a folder from the roots"`
		Sync     SyncCommandStruct   `subcommand:"sync,sync the roots"`
		List     ListCommandStruct   `subcommand:"list,list the roots"`
		Find     FindCommandStruct   `subcommand:"find,find a folder"`
		SetupCmd SetupCommandStruct  `subcommand:"setup,setup the current folder configuration"`
	}

	AddCommandStruct struct {
		baseDevStruct

		Args []string `args:"1,folder path to add\n(can be a relative, absolute and ~ resolved paths)"`
	}
	DeleteCommandStruct struct {
		baseDevStruct

		Args []string `args:"1,folder path to delete\n(can be a relative, absolute and ~ resolved paths)"`
	}
	SyncCommandStruct struct {
		baseDevStruct
	}
	ListCommandStruct struct {
		baseDevStruct
	}
	FindCommandStruct struct {
		baseDevStruct

		Args []string `args:"-1,query string to find a folder"`
	}
	SetupCommandStruct struct {
		baseDevStruct
	}
)

func (ds *DevStruct) Run(ctx context.Context, output io.Writer) error {
	return errors.NewError(nil, "a subcommand is required", false)
}

func (ac *AddCommandStruct) Run(ctx context.Context, output io.Writer) error {
	_, _ = fmt.Fprintln(output, "Adding folder to roots")

	root, err := ac.service.CanAddRoot(ac.Args[0])
	if err != nil {
		return err
	}

	return ac.service.AddRoot(root)
}

func (dc *DeleteCommandStruct) Run(ctx context.Context, output io.Writer) error {
	return dc.service.DeleteRoot(dc.Args[0])
}

func (sc *SyncCommandStruct) Run(ctx context.Context, output io.Writer) error {
	return sc.service.Sync()
}

func (lc *ListCommandStruct) Run(ctx context.Context, output io.Writer) error {
	return lc.service.ListRoots()
}

func (fc *FindCommandStruct) Run(ctx context.Context, output io.Writer) error {
	return fc.service.RunFind(fc.Args)
}

func (sc *SetupCommandStruct) Run(ctx context.Context, output io.Writer) error {
	_, _ = fmt.Fprintln(output, "Running setup command")
	return sc.service.Setup()
}

func (bs *baseDevStruct) Setup(ctx context.Context) error {
	cd, err := contextdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	bs.service = devservice.NewService(cd.RootConfig)

	return nil
}
