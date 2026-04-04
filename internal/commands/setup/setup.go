package setup

import (
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	configservice "github.com/guionardo/gs-dev/internal/services/config"
	"github.com/spf13/cobra"
)

type SetupCommand struct {
	commands.CommonCommand

	configService *configservice.ConfigService
}

const (
	name        = "setup"
	description = "Setup the application"
)

func (s *SetupCommand) Init() {
	s.InitCommandVariables(name, description, false, "", false)
}
func (s *SetupCommand) Setup(configuration *config.ConfigFile) error {
	s.configService = configservice.NewConfigService(configuration)
	return nil
}

func (s *SetupCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		RunE:  s.Run,
	}
}

func (s *SetupCommand) Run(cmd *cobra.Command, args []string) error {
	return nil
}

func (s *SetupCommand) GetTUICommand() func() error {
	return func() error {
		return s.configService.EditConfig()
	}
}
