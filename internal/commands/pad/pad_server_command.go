package padcommand

import (
	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"

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

	_, err := config.GetValue[padservice.PadServerConfig](configuration)

	return err
}

func (p *PadServerCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&PadServerStruct{}, serverName, serverDescription, "", false)
}

func (p *PadServerCommand) GetTUICommand() func() error {
	return func() error {
		return nil //TODO: Implement
	}
}
