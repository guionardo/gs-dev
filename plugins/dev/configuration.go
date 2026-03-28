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
		LastChosen      map[string]int `yaml:"last_chosen"`
		MostChosenCount int            `yaml:"most_chosen_count" default:"5"`
	}
	RootsConfiguration struct {
		Roots map[string]Root `yaml:"roots"`
	}

	LostChosenFolders map[string]int
)

func (c *Configuration) IncrementChosen(folder string) {
	c.LastChosen[folder]++
}
