package dev

import (
	"github.com/guionardo/gs-dev/configs"
	"gopkg.in/yaml.v3"
)

type DevConfig struct {
	configs.ConfigFolder `yaml:"-"`
	DevFolders           []string `yaml:"dev_folders"`
	MaxSubLevels         int      `yaml:"max_sub_levels"`
}

func (cfg *DevConfig) Load() error {
	if content, err := cfg.LoadFile("dev"); err != nil {
		return err
	} else {
		return yaml.Unmarshal(content, cfg)
	}
}

func (cfg *DevConfig) Save() error {
	return cfg.SaveFile("dev", cfg)
}

func NewDevConfig(root *configs.ConfigFolder) (cfg *DevConfig) {
	cfg = &DevConfig{
		ConfigFolder: *root,
		DevFolders:   []string{},
		MaxSubLevels: 0,
	}

	return cfg
}
