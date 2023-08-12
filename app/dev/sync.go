package dev

import (
	"time"

	"github.com/guionardo/gs-dev/internal/logger"
)

func RunSync(maxSyncInterval time.Duration) error {

	devConfig := getConfig()

	for _, rootFolder := range devConfig.Folders {
		if err := rootFolder.Sync(); err != nil {
			return err
		}
	}
	devConfig.LastSync = time.Now()
	if maxSyncInterval.Minutes() > 0 {
		devConfig.MaxSyncInterval = maxSyncInterval
		logger.Info("Updated max synchronization interval to %v", maxSyncInterval)
	}
	return devConfig.Save()
}
