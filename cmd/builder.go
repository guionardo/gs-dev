package cmd

import (
	"errors"
	"log/slog"

	"github.com/guionardo/gs-dev/internal/config"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/logging"
	"github.com/spf13/cobra"
)

type (
	command[TParams any, TConfig config.ConfigType] struct {
		name       string
		configRoot *config.ConfigRoot
		params     TParams
		config     TConfig
	}

	Commander interface {
		GetCobraCommand() *cobra.Command

		setConfig(any) error
	}
	CommandOptionsFn func(Commander) error
)

// NewCommand Create a command binded to a ConfigType
func NewCommand[TParams any, TConfig config.ConfigType](name string, configRoot *config.ConfigRoot, options ...CommandOptionsFn) (Commander, error) {
	if configRoot == nil {
		return nil, errs.NewError(nil, "null configuration", false)
	}

	config, err := config.GetValue[TConfig](configRoot)
	if err != nil {
		logging.Debug("Configuration", slog.Any("error", err), slog.Any("config", config))
		// TODO: Check if the configuration should be valid at this point
	}

	cmd := &command[TParams, TConfig]{
		name:       name,
		configRoot: configRoot,
		config:     config,
	}
	for _, option := range options {
		err := option(cmd)
		if err, ok := errors.AsType[errs.Error](err); ok {
			if !err.IsRecoverable() {
				return nil, err
			}
		}
	}

	return cmd, nil
}

func (c *command[TParams, TConfig]) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: c.name,
	}

	return cmd
}

func (c *command[TParams, TConfig]) setConfig(config any) error {
	if typedConfig, ok := config.(TConfig); ok {
		c.config = typedConfig
		return nil
	} else {
		return errs.NewError(nil, "expected configuration of type %T for command %s. Got %T", false, typedConfig, config)
	}
}
