package console

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func validateFolder(folderName string) error {
	if _, err := os.Stat(folderName); err != nil {
		return errors.New("path not found")
	}
	return nil
}

func PromptFolder(label string, defaultFolder string) (folder string, err error) {
	for {
		prompt := promptui.Prompt{Label: "Configuration path",
			Validate: validateFolder,
			Default:  defaultFolder,
		}
		folder, err := prompt.Run()

		if err == nil {
			return folder, nil
		}
		retryPrompt := promptui.Prompt{
			Label:     fmt.Sprintf("Error: %v. Retry?", err),
			IsConfirm: true,
			Default:   "Y"}
		if response, err := retryPrompt.Run(); err == nil {
			response = strings.ToLower(response + " ")
			if !strings.ContainsAny(fmt.Sprintf("%v", response[0]), "yst1") {
				return "", errors.New("path selection aborted")
			}
		}
	}
}
