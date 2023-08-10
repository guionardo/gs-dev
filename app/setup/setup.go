package setup

import (
	"errors"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/guionardo/gs-dev/app/dev"
	"github.com/spf13/cobra"
)

var (
	rootCmd *cobra.Command
)

func SetRootCmd(root *cobra.Command) {
	rootCmd = root
}

// Parse setup from commands
func Setup(cmd *cobra.Command) {
	mainCommands := make([]string, 1, 10)
	for _, command := range rootCmd.Commands() {
		if _, ok := dev.InteractiveCommands[command.Use]; ok {
			mainCommands = append(mainCommands, command.Use)
		}
	}
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
			function := dev.InteractiveCommands[answers.MainCommand]
			err = function(cmd)
		}
	}

	if err != nil {
		cmd.PrintErr(err)
		return
	}

	// // the questions to ask
	// var qs = []*survey.Question{
	// 	{
	// 		Name:      "name",
	// 		Prompt:    &survey.Input{Message: "What is your name?"},
	// 		Validate:  survey.Required,
	// 		Transform: survey.Title,
	// 	},
	// 	{
	// 		Name: "color",
	// 		Prompt: &survey.Select{
	// 			Message: "Choose a color:",
	// 			Options: []string{"red", "blue", "green"},
	// 			Default: "red",
	// 		},
	// 	},
	// 	{
	// 		Name:   "age",
	// 		Prompt: &survey.Input{Message: "How old are you?"},
	// 	},
	// }
	// // the answers will be written to this struct
	// answers := struct {
	// 	Name          string // survey will match the question and field names
	// 	FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
	// 	Age           int    // if the types don't match, survey will convert it
	// }{}

	// perform the questions
	// err := survey.Ask(qs, &answers)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }

	fmt.Printf("chose %s.", answers.MainCommand)
	fmt.Scanln()

}
