package configs

import (
	"fmt"
	"os"
	"strconv"
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

func PromptInt(message string, defaultValue int, minValue int, maxValue int) (value int, err error) {
	prompt := promptui.Prompt{
		Label:   fmt.Sprintf("%s (default=%d)", message, defaultValue),
		Default: fmt.Sprintf("%d", defaultValue),
		Validate: func(s string) error {
			i, err := strconv.Atoi(s)
			if err != nil {
				return fmt.Errorf("invalid number %s", s)
			}
			if i < minValue || i > maxValue {
				return fmt.Errorf("%s must be between %d and %d", message, minValue, maxValue)
			}
			return nil
		},
	}
	strValue, err := prompt.Run()
	if err == nil {
		value, err = strconv.Atoi(strValue)
	}
	return
}

func PromptOpt(message string, items []string) (index int, value string, err error) {
	prompt := promptui.Select{
		Label: message,
		Items: items,
	}
	index, result, err := prompt.Run()
	return index, result, err
}
