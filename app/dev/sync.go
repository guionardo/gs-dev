package dev

import (
	"time"

	"github.com/guionardo/gs-dev/config"
)

func RunSync(cfg ...*config.DevConfig) error {
	var devConfig *config.DevConfig
	if len(cfg) > 0 {
		devConfig = cfg[0]
	} else {
		devConfig = getConfig()
	}
	for _, rootFolder := range devConfig.Folders {
		if err := rootFolder.Sync(); err != nil {
			return err
		}
	}
	devConfig.LastSync = time.Now()
	return devConfig.Save()
}
