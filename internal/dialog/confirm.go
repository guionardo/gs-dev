package dialog

import "github.com/AlecAivazis/survey/v2"

func Confirm(label string, defaultValue bool) bool {
	confirm := &survey.Confirm{
		Message: label,
		Default: defaultValue,
	}

	confirmed := false
	if err := survey.AskOne(confirm, &confirmed); err != nil {
		return false
	}

	return confirmed
}
