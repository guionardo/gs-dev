package padcommand

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	ctxdata "github.com/guionardo/gs-dev/internal/context"
	"github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/logging"
	padserver "github.com/guionardo/gs-dev/internal/pad/server"
	"github.com/guionardo/gs-dev/internal/pad/storage"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"
	"github.com/guionardo/gs-dev/pkg/console"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

type (
	PadServerStruct struct {
		padServerStruct

		Serve   PadServerServeStruct `subcommand:"serve"`
		DoSetup PadServerSetupStruct `subcommand:"setup"`
	}

	PadServerServeStruct struct {
		padServerStruct

		Port              int           `flag:"port" description:"listening port" default:"0"`
		APIKey            string        `flag:"api-key" description:"API key for authorization" default:""`
		StorageDirectory  string        `flag:"storage" description:"storage directory" default:""`
		StorageDefaultTTL time.Duration `flag:"ttl" description:"default TTL" default:"24h"`
	}

	PadServerSetupStruct struct {
		padServerStruct

		Show              bool          `flag:"show" description:"Show current server setup"`
		Port              int           `flag:"port" description:"listening port"`
		APIKey            string        `flag:"api-key" description:"API key for authorization"`
		StorageDirectory  string        `flag:"storage" description:"storage directory"`
		StorageDefaultTTL time.Duration `flag:"ttl" description:"default TTL" default:"24h"`
	}

	padServerStruct struct {
		configRoot *config.ConfigRoot
	}
)

func (p *padServerStruct) Setup(ctx context.Context) error {
	cd, err := ctxdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	p.configRoot = cd.RootConfig

	return nil
}

func (p *padServerStruct) getConfig() padservice.PadServerConfig {
	cfg, _ := config.GetValue[padservice.PadServerConfig](p.configRoot)
	return cfg
}

func (p *PadServerStruct) Run(ctx context.Context, output io.Writer) error {
	return errors.NewError(nil, "a subcommand is required", false)
}

func (p *PadServerServeStruct) Run(ctx context.Context, output io.Writer) error {
	cfg := p.getConfig()
	if err := cfg.Validate(); err != nil {
		return err
	}

	flagsSet := cli.GetFlagsSet(ctx)

	if flagsSet.Has("storage") {
		cfg.StorageConfig.Directory = p.StorageDirectory
	}

	if flagsSet.Has("ttl") {
		cfg.StorageConfig.DefaultTTL = p.StorageDefaultTTL
	}

	if flagsSet.Has("port") {
		cfg.Port = p.Port
	}

	if flagsSet.Has("api-key") {
		cfg.APIKey = p.APIKey
	}

	storage, err := storage.NewFileSystemStorage(cfg.StorageConfig, ctx, logging.Logger())
	if err != nil {
		return err
	}

	service, err := padservice.NewPadServerService(p.configRoot, storage)
	if err != nil {
		return err
	}

	server := padserver.NewPadGrpcServer(&cfg, service, logging.Logger())

	return server.Start(ctx)
}

func (p *PadServerSetupStruct) Run(ctx context.Context, output io.Writer) error {
	cfg := p.getConfig()
	cfg.Defaults()

	if p.Show {
		return p.showSetup(cfg, output)
	}

	flags := cli.GetFlagsSet(ctx)
	if flags.Has("storage") {
		cfg.StorageConfig.Directory = p.StorageDirectory
	}

	if flags.Has("ttl") {
		cfg.StorageConfig.DefaultTTL = p.StorageDefaultTTL
	}

	if flags.Has("port") {
		cfg.Port = p.Port
	}

	if flags.Has("api-key") {
		cfg.APIKey = p.APIKey
	}

	if cfg.StorageConfig.Directory != "" {
		sd, err := files.AssertDirectory(cfg.StorageConfig.Directory)
		if err != nil {
			if err = os.MkdirAll(cfg.StorageConfig.Directory, consts.DirPermissions); err == nil {
				sd, err = files.AssertDirectory(cfg.StorageConfig.Directory)
			}
		}

		if err == nil {
			cfg.StorageConfig.Directory = sd
		}
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	config.SetValue(p.configRoot, cfg)

	if err := p.configRoot.Save(); err != nil {
		return err
	}

	return p.showSetup(cfg, output)
}

func (p *PadServerSetupStruct) showSetup(cfg padservice.PadServerConfig, output io.Writer) error {
	tree := console.NewTree(console.Styled("Pad server configuration", console.Blue))
	tree.AddChild(console.Styled("PORT = %d", console.Cyan, cfg.Port))
	tree.AddChild(console.Styled("API KEY = %s", console.Cyan, cfg.APIKey))
	tree.AddChild(console.Styled("STORAGE = %s", console.Cyan, cfg.StorageConfig.Directory))
	tree.AddChild(console.Styled("S = %s", console.Cyan, cfg.StorageConfig.DefaultTTL))

	if err := cfg.Validate(); err == nil {
		tree.AddChild(console.Styled("Status: OK", console.Green))
	} else {
		tree.AddChild(console.Styled("Status: %s", console.Red, err.Error()))
	}

	return tree.Write(output)
}
