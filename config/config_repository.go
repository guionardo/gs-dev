package config

import pathtools "github.com/guionardo/gs-dev/internal/path_tools"

var configRepoFolder = ""

func ValidateRepositoryFolder(folder string) error {
	if err := pathtools.ValidatePath(folder); err == nil {
		configRepoFolder = folder
		return nil
	} else {
		return err
	}
}

func GetConfigRepositoryFolder() string {
	return configRepoFolder
}
