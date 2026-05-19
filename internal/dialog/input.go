package dialog

import (
	"github.com/AlecAivazis/survey/v2"
	errs "github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/pkg/tools/files"
)

func InputDirectory(label string, validate func(string) error) (string, error) {
	input := &survey.Input{
		Message: label,
		Default: ".",
	}

	directory := ""
	if err := survey.AskOne(input, &directory, survey.WithValidator(func(ans any) error {
		if validate != nil {
			return validate(ans.(string))
		}

		return nil
	})); err != nil {
		return "", errs.NewError(err, "error validating directory", false)
	}

	return files.AssertDirectory(directory)
}

func Input(label string) (string, error) {
	input := &survey.Input{
		Message: label,
	}

	value := ""
	if err := survey.AskOne(input, &value); err != nil {
		return "", errs.NewError(err, "error inputing value", false)
	}

	return value, nil
}
