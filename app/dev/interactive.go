package dev

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/logger"
	"github.com/spf13/cobra"
)

var mainCommands = make([]string, 0) //TODO: Implementar
type SetupFunction func(cmd *cobra.Command) error

var InteractiveCommands = make(map[string]SetupFunction)

func SetInteractiveCommand(command *cobra.Command, setupFunction SetupFunction) {
	InteractiveCommands[command.Use] = setupFunction
}

func DevInteractive(cmd *cobra.Command) error {
	cfg := configs.NewConfigFolder(cmd)
	if cfg.ErrorCode > 0 {
		return cfg.Error
	}

	// Check
	var questions = []*survey.Question{
		{
			Name: "main_command",
			Prompt: &survey.Select{
				Message: "Choose a command:",
				Options: mainCommands,
			},
		}}

	answers := struct {
		MainCommand string `survey:"main_command"`
	}{}

	var err error
	if err = survey.Ask(questions, &answers); err == nil {
		if len(answers.MainCommand) == 0 {
			err = errors.New("command is required")
		} else {
			function := InteractiveCommands[answers.MainCommand]
			err = function(cmd)
		}
	}
	logger.Debug("DEV")

	return err
}
