package dev

import (
	"errors"
	"fmt"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/guionardo/gs-dev/internal"
)

func RunGo(args []string, output string) error {
	os.Remove(output)
	if subFolders := findArgs(args); len(subFolders) == 0 {
		return fmt.Errorf("there is no folder for args %v", args)
	} else {
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

		if err := survey.Ask(questions, &answers); err != nil {
			return err
		}

		if len(answers.Folder) == 0 {
			return errors.New("folder is required")
		}
		os.WriteFile(output, []byte("cd "+answers.Folder), 0666)
	}
	return nil
}

func findArgs(args []string) (folders []string) {
	foundFolders := make(map[string]struct{})
	devConfig := getConfig()
	for _, folder := range devConfig.Folders {
		for _, path := range folder.FindByPattern(args) {
			foundFolders[path.Name] = struct{}{}
		}

	}
	return internal.List(foundFolders)
}
