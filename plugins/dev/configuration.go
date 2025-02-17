package dev

import (
	"time"

	"github.com/guionardo/gs-dev/plugins/commons"
)

type (
	Root struct {
		Folders  []string `yaml:"folders"`
		MaxDepth int      `yaml:"max_depth"`
	}
	Configuration struct {
		commons.PluginConfiguration
		LastSync        time.Time      `yaml:"last_sync"`
		SyncInterval    time.Duration  `yaml:"sync_interval"`
		DefaultMaxDepth int            `yaml:"default_max_depth"`
		LastChoosen     map[string]int `yaml:"last_choosen"`
	}
	RootsConfiguration struct {
		Roots map[string]Root `yaml:"roots"`
	}
	LocalConfig struct {
		Ignore           bool   `yaml:"ignore"`
		IgnoreSubfolders bool   `yaml:"ignore_subfolders"`
		Description      string `yaml:"description"`
	}
	LastChoosenFolders map[string]int
)

func (c *Configuration) IncrementChoosed(folder string) {
	c.LastChoosen[folder]++
}
