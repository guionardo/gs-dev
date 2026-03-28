package devservice

import (
	"time"

	"github.com/guionardo/gs-dev/plugins/commons"
)

type (
	DevConfiguration struct {
		commons.PluginConfiguration

		LastSync        time.Time      `yaml:"last_sync"`
		SyncInterval    time.Duration  `yaml:"sync_interval"`
		DefaultMaxDepth int            `yaml:"default_max_depth"`
		LastChosen      map[string]int `yaml:"last_chosen"`
		MostChosenCount int            `yaml:"most_chosen_count" default:"5"`
	}
	Root struct {
		Folders  []string `yaml:"folders"`
		MaxDepth int      `yaml:"max_depth"`
	}
	// RootsConfiguration is a map[rootPath]Root
	RootsConfiguration map[string]Root

	LostChosenFolders map[string]int
)

const (
	DefaultMostChosenCount = 5
	DefaultMaxDepth        = 3
	DefaultSyncInterval    = time.Hour
)

func (c *DevConfiguration) Defaults() {
	if c.MostChosenCount < 1 {
		c.MostChosenCount = DefaultMostChosenCount
	}
	if c.DefaultMaxDepth < 1 {
		c.DefaultMaxDepth = DefaultMaxDepth
	}
	if c.SyncInterval == 0 {
		c.SyncInterval = DefaultSyncInterval
	}
	if c.LastChosen == nil {
		c.LastChosen = make(map[string]int)
	}
}

func (c *Root) Defaults() {
	if c.MaxDepth < 1 {
		c.MaxDepth = DefaultMaxDepth
	}
	if c.Folders == nil {
		c.Folders = make([]string, 0)
	}
}
