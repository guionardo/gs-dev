package configservice

import (
	"github.com/guionardo/gs-dev/internal/colors"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	devservice "github.com/guionardo/gs-dev/internal/services/dev"
)

type ConfigService struct {
	configFile *config.ConfigFile
}

func NewConfigService(configFile *config.ConfigFile) *ConfigService {
	return &ConfigService{
		configFile: configFile,
	}
}

func (c *ConfigService) EditConfig() error {
	devConfig, err := config.GetValue[devservice.DevConfiguration](c.configFile, "dev")
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

	config.SetValue(c.configFile, "dev", devConfig)

	if err = config.SetValue(c.configFile, "dev", devConfig); err != nil {
		colors.Error("Error saving config: %v", err)
		return err
	}

	colors.Success("Config saved successfully")

	return nil
}
