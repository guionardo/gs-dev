package padcommand

import (
	"errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/configurations"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/logging"
	padserver "github.com/guionardo/gs-dev/internal/pad/server"

	// padservice "github.com/guionardo/gs-dev/internal/pad/service"
	"github.com/guionardo/gs-dev/internal/pad/storage"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"
	"github.com/spf13/cobra"
)

type PadServerCommand struct {
	commands.CommonCommand

	configuration *config.ConfigRoot
}

const (
	serverName        = "pad-server"
	serverDescription = "Start the Pad server."
)

func (p *PadServerCommand) Init() {
	p.InitCommandVariables(serverName, serverDescription, p, p.setup)
}

func (p *PadServerCommand) setup(configuration *config.ConfigRoot) error {
	p.configuration = configuration

	_, err := config.GetValue[configurations.PadServerConfig](configuration)

	return err
}

func (p *PadServerCommand) BuildCommand() *cobra.Command {
	serverConfig, _ := config.GetValue[configurations.PadServerConfig](p.configuration)
	cmd := &cobra.Command{
		Use:   serverName,
		Short: serverDescription,
		RunE:  p.runCommandServer,
	}

	cmd.Flags().IntP("port", "p", serverConfig.Port, "Port to run the Pad server on")
	cmd.Flags().String("api-key", "", "API key for authenticating with the Pad server (overrides config)")

	return cmd
}

func (p *PadServerCommand) GetTUICommand() func() error {
	return func() error {
		return nil //TODO: Implement
	}
}

func (p *PadServerCommand) runCommandServer(cmd *cobra.Command, args []string) error {
	ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	store, err := storage.NewStorage(p.configuration, ctx)
	if err != nil {
		return err
	}

	padServerService, err := padservice.NewPadServerService(p.configuration, store)
	if err != nil {
		return err
	}

	serverConfig := padServerService.GetConfig()

	if cmd.Flags().Changed("port") {
		port, _ := cmd.Flags().GetInt("port")
		serverConfig.Port = port
	}

	if cmd.Flags().Changed("api-key") {
		apiKey, _ := cmd.Flags().GetString("api-key")
		serverConfig.APIKey = apiKey
	}

	if err := serverConfig.Validate(); err != nil {
		if e, ok := errors.AsType[errs.Error](err); ok {
			if e.IsRecoverable() {
				serverConfig.Defaults()
			}
		}

		return err
	}

	server := padserver.NewPadGrpcServer(serverConfig, padServerService, logging.Logger())

	return server.Start(ctx)
}
