package dev

import (
	"time"
)

func RunSync() error {
	devConfig := getConfig()
	for _, rootFolder := range devConfig.Folders {
		if err := rootFolder.Sync(); err != nil {
			return err
		}
	}
	devConfig.LastSync = time.Now()
	return devConfig.Save()
}
