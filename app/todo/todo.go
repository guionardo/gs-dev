package todo

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

const NewTag = "* NEW"

func RunAddTodo(cmd *cobra.Command, args []string) error {
	config := getConfig()
	tags := append([]string{NewTag}, config.GetAllTags()...)
	questions := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "TO-DO"},
			Validate: survey.Required,
		},
		{
			Name: "due-to",
			Prompt: &survey.Select{
				Message: "Due to",
				Options: []string{
					"Infinity",
					"Tomorrow",
					"2 days",
					"One week",
					"Custom date",
				},
			},
			Validate: survey.Required,
		},
		{
			Name: "tags",
			Prompt: &survey.MultiSelect{
				Message: "Tags",
				Options: tags,
			},
		},
	}
	answers := struct {
		Name  string
		DueTo string   `survey:"due-to"`
		Tags  []string `survey:"tags"`
	}{}

	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	return nil
}
