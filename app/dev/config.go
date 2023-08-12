package dev

import "github.com/guionardo/gs-dev/config"

func getConfig() *config.DevConfig {
	devConfig := config.NewDevConfig(config.GetConfigRepositoryFolder())
	_ = devConfig.Load()
	return devConfig
}
