package dev

import (
	"fmt"
	"strings"

	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

func RunRemoveFolder(folder string) error {
	folder = pathtools.AbsPath(folder)
	devConfig := getConfig()
	if _, ok := devConfig.Folders[folder]; !ok {
		return fmt.Errorf("path was not previously defined - %s", folder)
	}

	// Removes all the subfolders
	toRemove := make([]string, 0, len(devConfig.Folders))
	for f := range devConfig.Folders {
		if strings.HasPrefix(f, folder) {
			toRemove = append(toRemove, f)
		}
	}
	for _, f := range toRemove {
		delete(devConfig.Folders, f)
		fmt.Printf("Removed %s\n", f)
	}

	if err := devConfig.Save(); err != nil {
		return err
	}

	return nil
}
