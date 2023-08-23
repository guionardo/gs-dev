package dev

import (
	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/internal/logger"
)

var readenConfig *config.DevConfig

func getConfig() *config.DevConfig {
	if readenConfig == nil {
		readenConfig = config.NewDevConfig(config.GetConfigRepositoryFolder())
		if err := readenConfig.Load(); err != nil {
			logger.Error("Failed to read DevConfig %v", err)
		} else {
			logger.Debug("Loaded DevConfig")
		}
	}

	return readenConfig
}
