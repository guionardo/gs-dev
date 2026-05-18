package padcommand

import (
	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/spf13/cobra"
)

type PadCliCommand struct {
	commands.CommonCommand
}

const (
	cliName        = "pad"
	cliDescription = "Pad is a form of transmit short information."
)

func (p *PadCliCommand) Init() {
	p.InitCommandVariables(cliName, cliDescription, p, p.setup).WithInitAlias().WithOutput()
}

func (p *PadCliCommand) setup(configuration *config.ConfigRoot) (err error) {
	return nil
}

func (p *PadCliCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&PadCliStruct{}, cliName, cliDescription, "", false)
}

func (p *PadCliCommand) GetTUICommand() func() error {
	return func() error {
		return nil //TODO: Implement
	}
}
