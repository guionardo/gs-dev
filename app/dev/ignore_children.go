package dev

import (
	"fmt"

	"github.com/guionardo/go-dev/pkg/logger"
	"github.com/guionardo/gs-dev/internal"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

func RunIgnoreChildren(folderName string, ignoreChildren bool) error {
	folderName = pathtools.AbsPath(folderName)
	devConfig := getConfig()
	folder := devConfig.GetPath(folderName)
	if folder == nil {
		return fmt.Errorf("path was not previously defined - %s", folderName)
	}
	log := ""
	if ignoreChildren {
		if folder.Status == internal.Disabled {
			return fmt.Errorf("path is disabled already - %s", folderName)
		}
		if folder.Status == internal.ChildrenDisabled {
			return fmt.Errorf("path is already ignoring children - %s", folderName)
		}
		folder.Status = internal.ChildrenDisabled
		log = "Disabling"
	} else {
		if folder.Status == internal.Enabled {
			return fmt.Errorf("path is already accepting children %s", folderName)
		}
		folder.Status = internal.Enabled
		log = "Enabling"
	}
	if err := devConfig.Save(); err != nil {
		return err
	}
	if err := RunSync(); err == nil {
		logger.Info("%s children of %s", log, folderName)
		return nil
	} else {
		return err
	}

}
