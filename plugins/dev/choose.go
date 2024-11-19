package dev

import (
	"errors"

	"github.com/AlecAivazis/survey/v2"
)

func chooseFolder(subFolders []string) (folder string, err error) {
	if len(subFolders) == 1 {
		return subFolders[0], nil
	}

	var questions = []*survey.Question{
		{
			Name: "folder",
			Prompt: &survey.Select{
				Message: "Choose a folder:",
				Options: subFolders,
			},
		}}

	answers := struct {
		Folder string `survey:"folder"`
	}{}

	if err = survey.Ask(questions, &answers); err == nil {
		if len(answers.Folder) == 0 {
			err = errors.New("folder is required")
		} else {
			folder = answers.Folder
		}
	}
	return
}
