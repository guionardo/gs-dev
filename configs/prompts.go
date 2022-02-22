package configs

import (
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

func Confirm(message string, defaultValue bool) bool {
	dV := "y"
	if !defaultValue {
		dV = "n"
	}
	prompt := promptui.Prompt{
		Label:     message,
		IsConfirm: true,
		Default:   dV,
	}
	result, err := prompt.Run()
	return err == nil && strings.HasPrefix(strings.ToUpper(result), "Y")

}

func Prompt(message string, defaultValue string) (value string, err error) {
	prompt := promptui.Prompt{
		Label:   message,
		Default: defaultValue,
	}
	return prompt.Run()
}

func PromptPath(message string, defaultPath string) (value string, err error) {
	validatePath := func(path string) error {
		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("path not found [%s]", path)
		}
		return nil
	}
	prompt := promptui.Prompt{
		Label:    message,
		Default:  defaultPath,
		Validate: validatePath,
	}
	return prompt.Run()
}
