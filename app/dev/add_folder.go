package dev

import (
	"fmt"

	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/internal"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

func RunAddFolder(folder string, max_sublevels uint8) error {
	folder = pathtools.AbsPath(folder)
	if !pathtools.PathExists(folder) {
		return fmt.Errorf("path not found - %s", folder)
	}
	devConfig := getConfig()

	if p := findSubFolder(devConfig, folder); p != nil {
		return fmt.Errorf("path was previously defined - %s", p.Name)
	}

	path := internal.NewPathList(folder, int(max_sublevels))
	if err := path.Sync(); err != nil {
		return err
	}
	devConfig.Folders[folder] = path

	if err := devConfig.Save(); err != nil {
		return err
	}
	fmt.Printf("Folder added %s\n", folder)
	return nil
}

func findSubFolder(devConfig *config.DevConfig, folder string) *internal.Path {
	for _, rootFolder := range devConfig.Folders {
		if p := rootFolder.Find(folder); p != nil {
			return p
		}
	}
	return nil
}
