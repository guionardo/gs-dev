package devservice

import (
	"path"
	"sort"
	"time"

	pathtools "github.com/guionardo/go/path_tools"
)

type (
	DevConfiguration struct {
		Enabled         bool           `yaml:"enabled" default:"true"`
		LastSync        time.Time      `yaml:"last_sync"`
		SyncInterval    time.Duration  `yaml:"sync_interval" default:"1h"`
		DefaultMaxDepth int            `yaml:"default_max_depth" default:"3"`
		LastChosen      map[string]int `yaml:"last_chosen"`
		MostChosenCount int            `yaml:"most_chosen_count" default:"5"`
	}

	Root struct {
		Folders      []string               `yaml:"folders"`
		MaxDepth     int                    `yaml:"max_depth"`
		LocalConfigs map[string]LocalConfig `yaml:"local_configs"`
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

func (c DevConfiguration) Key() string {
	return "dev"
}

func (c DevConfiguration) Validate() error {
	return nil // TODO: Implementar validação
}

func (c *Root) Validate() error {
	return nil
}

func (c *Root) Defaults() {
	if c.MaxDepth < 1 {
		c.MaxDepth = DefaultMaxDepth
	}

	if c.Folders == nil {
		c.Folders = make([]string, 0)
	}

	if c.LocalConfigs == nil {
		c.LocalConfigs = make(map[string]LocalConfig)
	}
}

func (r *Root) EfectiveIgnore(folder string) bool {
	currentFolder := folder
	level := 0

	for !pathtools.IsRootDirectory(currentFolder) {
		if localConfig, ok := r.LocalConfigs[folder]; ok && (localConfig.Ignore || (level > 0 && localConfig.IgnoreSubfolders)) {
			return true
		}

		currentFolder = path.Dir(currentFolder)
		level++
	}

	return false
}

func (r *Root) CanIncludeFolder(folder string) bool {
	return !r.EfectiveIgnore(folder)
}

// Resync syncs the root with the local configs
func (r *Root) Resync() (changed bool) {
	for folder, rootLocalConfig := range r.LocalConfigs {
		if r.EfectiveIgnore(folder) {
			continue
		}

		localConfig, err := NewLocalConfig(folder)
		if err != nil {
			delete(r.LocalConfigs, folder)

			changed = true

			continue
		}

		if rootLocalConfig.Equal(*localConfig) {
			continue
		}

		changed = true
		r.LocalConfigs[folder] = *localConfig
	}

	if changed {
		r.Folders = make([]string, 0)
		for folder := range r.LocalConfigs {
			r.Folders = append(r.Folders, folder)
		}

		sort.Strings(r.Folders)
	}

	return changed
}

func (rc RootsConfiguration) Key() string {
	return "roots"
}

func (rc RootsConfiguration) Validate() error {
	//TODO: Implementar validação
	return nil
}
