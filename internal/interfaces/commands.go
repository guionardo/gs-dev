package interfaces

import (
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/spf13/cobra"
)

type Command interface {
	// Initialize the command
	Init()
	Setup(configuration *config.ConfigRoot) error
	GetName() string
	GetDescription() string

	// if the command is a cobra command, it will be added to the root command
	GetCobraCommand() *cobra.Command

	// if the command is a TUI command, it will be added to the TUI command
	GetTUICommand() func() error

	// if the command has an init alias, it will be added to the init script
	HasInitAlias() bool

	// Use by the init alias
	GetArguments() string

	UsesOutput() bool
}
