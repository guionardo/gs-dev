package dev

import "github.com/guionardo/gs-dev/config"

var readenConfig *config.DevConfig

func getConfig() *config.DevConfig {
	if readenConfig == nil {
		readenConfig = config.NewDevConfig(config.GetConfigRepositoryFolder())
		_ = readenConfig.Load()
	}

	return readenConfig
}
