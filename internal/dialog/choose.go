package dialog

import (
	"errors"
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
)

type ChooseItem struct {
	Name        string
	Description string
}

func NewChooseItem(name string, description string) ChooseItem {
	return ChooseItem{
		Name:        name,
		Description: description,
	}
}

// Choose a option from a list of options
// Returns the selected option and an error if the selection is invalid
// If there is only one option, it returns the option and nil
func Choose(label string, options ...ChooseItem) (answer string, err error) {
	switch len(options) {
	case 0:
		return "", errors.New("no options provided")
	case 1:
		return options[0].Name, nil
	}

	sort.Slice(options, func(i, j int) bool {
		return options[i].Name < options[j].Name
	})

	values := make([]string, len(options))

	for index, item := range options {
		values[index] = item.Name
	}

	var questions = []*survey.Question{
		{
			Name: "option",
			Prompt: &survey.Select{
				Message: label,
				Options: values,
				Description: func(value string, index int) string {
					return options[index].Description
				},
			},
		}}

	answers := struct {
		Option string `survey:"option"`
	}{}

	if err = survey.Ask(questions, &answers); err == nil {
		if len(answers.Option) == 0 {
			return "", fmt.Errorf("%s is required", label)
		}

		return answers.Option, nil
	}

	return "", err
}

func ToAnyArray(items []string) []ChooseItem {
	chooseItems := make([]ChooseItem, len(items))
	for index, option := range items {
		chooseItems[index] = ChooseItem{
			Name: option,
		}
	}

	return chooseItems
}
