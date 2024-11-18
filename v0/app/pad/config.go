package pad

import (
	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/internal/logger"
)

var readenConfig *config.PadConfig

func getConfig() *config.PadConfig {
	if readenConfig == nil {
		readenConfig = config.NewPadConfig(config.GetConfigRepositoryFolder())
		if err := readenConfig.Load(); err != nil {
			logger.Error("Failed to read PadConfig %v - using an empty one", err)
			readenConfig = &config.PadConfig{}
		} else {
			logger.Debug("Loaded PadConfig")
		}
	}

	return readenConfig
}
