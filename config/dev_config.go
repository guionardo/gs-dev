package config

import (
	"path"
	"time"

	"github.com/guionardo/gs-dev/internal"
)

//go:generate go run ../gen/configs.go

type DevConfig struct {
	fileName        string
	Folders         map[string]*internal.PathList `yaml:"folders"`
	LastSync        time.Time                     `yaml:"last_sync"`
	MaxSyncInterval time.Duration                 `yaml:"max_sync_interval"`
}

func (dc *DevConfig) GetPath(path string) *internal.Path {
	for _, folder := range dc.Folders {
		if p := folder.Find(path); p != nil {
			return p
		}
	}
	return nil
}

type DevFolder struct {
	IgnoreChildren bool           `yaml:"ignore_children"`
	Ignore         bool           `yaml:"ignore"`
	MaxSubLevels   uint8          `yaml:"max_sub_levels"`
	Root           string         `yaml:"root"`
	SubFolders     []DevSubFolder `yaml:"sub_folders"`
}

type DevSubFolder struct {
	Name           string `yaml:"name"`
	Ignore         bool   `yaml:"ignore"`
	IgnoreChildren bool   `yaml:"ignore_children"`
}

func NewDevConfig(root string) (cfg *DevConfig) {
	return &DevConfig{
		fileName:        path.Join(root, "dev.yaml"),
		Folders:         make(map[string]*internal.PathList),
		MaxSyncInterval: time.Hour,
	}
}

func (cfg *DevConfig) ShoudResync() bool {
	return cfg.LastSync.Before(time.Now().Add(-cfg.MaxSyncInterval))
}
