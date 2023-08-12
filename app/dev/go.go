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
		folder, err := chooseFolder(subFolders)
		if err == nil {
			err = os.WriteFile(output, []byte("cd "+folder), 0666)
		}
		return err
	}
}

func findArgs(args []string) (folders []string) {
	foundFolders := make(map[string]struct{})
	devConfig := getConfig()
	if devConfig.ShoudResync() {
		_ = RunSync(devConfig)
	}
	for _, folder := range devConfig.Folders {
		for _, path := range folder.FindByPattern(args) {
			foundFolders[path.Name] = struct{}{}
		}

	}
	return internal.List(foundFolders)
}

func chooseFolder(subFolders []string) (folder string, err error) {
	if len(subFolders) == 1 {
		folder = subFolders[0]
		return
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
