package sample

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/spf13/cobra"
)

type (
	SampleCommand struct {
		commands.Command
	}
	SampleCommand2 struct {
		commands.Command
	}
)

func (s *SampleCommand) GetName() string {
	return "sample"
}

func (s *SampleCommand) GetDescription() string {
	return "Sample command"
}

func (s *SampleCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sample",
		Short: "Sample command",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Sample command Cobra")
			return nil
		},
	}
}

func (s *SampleCommand) GetTUICommand() func() error {
	return func() error {
		fmt.Println("Sample command TUI")
		return nil
	}
}

func (s *SampleCommand2) GetName() string {
	return "sample2"
}

func (s *SampleCommand2) GetDescription() string {
	return "Sample command 2"
}

func (s *SampleCommand2) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "sample2",
		Short: "Sample command 2",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Sample command 2 Cobra")
			return nil
		},
	}
}

func (s *SampleCommand2) GetTUICommand() func() error {
	return func() error {
		fmt.Println("Sample command 2 TUI")
		return nil
	}
}
