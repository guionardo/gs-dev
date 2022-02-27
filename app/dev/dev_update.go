package dev

import (
	"fmt"
	"time"

	"github.com/guionardo/gs-dev/configs"
)

func UpdatePaths(cfg *configs.RootConfig) (changes int, err error) {
	changes = 0
	existingPaths := make([]configs.DevPathConfig, 0)
	for _, devPathConfig := range cfg.DevConfig.Paths {
		if devPathConfig.Exists() {
			existingPaths = append(existingPaths, devPathConfig)
		} else {
			fmt.Printf("Removed path %s\n", devPathConfig.FullPath)
			changes++
		}
	}

	cfg.DevConfig.Paths = existingPaths
	for _, devFolder := range cfg.DevConfig.DevFolders {
		for _, existingFolder := range ReadFolders(devFolder, cfg.DevConfig.MaxSubLevels) {
			index, _ := cfg.DevConfig.FindPath(existingFolder)
			if index == -1 {
				cfg.DevConfig.Paths = append(cfg.DevConfig.Paths, configs.DevPathConfig{
					FullPath:       existingFolder,
					IgnoreSubPaths: false,
					IsHidden:       false,
				})
				fmt.Printf("Added new path %s\n", existingFolder)
				changes++
			}
		}
	}
	cfg.DevConfig.LastUpdate = time.Now()
	return
}
