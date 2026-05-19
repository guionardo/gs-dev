package configservice

import (
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
	"github.com/guionardo/gs-dev/pkg/console"
)

type ConfigService struct {
	configFile *config.ConfigRoot
}

func NewConfigService(configFile *config.ConfigRoot) *ConfigService {
	return &ConfigService{
		configFile: configFile,
	}
}

func (c *ConfigService) EditConfig() error {
	devConfig, err := config.GetValue[devservice.DevConfiguration](c.configFile)
	if err != nil {
		return err
	}

	cfg, err := dialog.ReadConfig(devConfig.SyncInterval, devConfig.DefaultMaxDepth, devConfig.MostChosenCount)
	if err != nil {
		return err
	}

	devConfig.SyncInterval = cfg.SyncInterval
	devConfig.DefaultMaxDepth = cfg.DefaultMaxDepth
	devConfig.MostChosenCount = cfg.MostChosenCount

	config.SetValue(c.configFile, devConfig)

	if err = c.configFile.Save(); err != nil {
		console.Error("Error saving config: %v", err)
		return err
	}

	console.Success("Config saved successfully")

	return nil
}
