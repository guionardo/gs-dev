package configs

import (
	"os"
	"strings"
	"time"
)

type (
	DevConfig struct {
		DevFolders     []string        `yaml:"dev_folders"`
		MaxSubLevels   int             `yaml:"max_sub_levels"`
		Paths          []DevPathConfig `yaml:"paths"`
		LastUpdate     time.Time       `yaml:"last_update"`
		UpdateInterval time.Duration   `yaml:"update_interval"`
	}
	DevPathConfig struct {
		FullPath       string   `yaml:"fp"`
		IgnoreSubPaths bool     `yaml:"ign_sub"`
		IsHidden       bool     `yaml:"hid"`
		AfterCommands  []string `yaml:"aft_cmd"`
	}
)

func (devPathConfig *DevPathConfig) Exists() bool {
	stat, err := os.Stat(devPathConfig.FullPath)
	return err == nil && stat.IsDir()
}

func (cfg *DevConfig) Prompt() bool {
	if !Confirm("Setup dev config?", true) {
		return false
	}
	return true
}

func (cfg *DevConfig) Find(args []string) []DevPathConfig {
	results := make([]DevPathConfig, 0)

	for _, path := range cfg.Paths {
		if !path.IsHidden && path.Match(args) {
			results = append(results, path)
		}
	}
	return results
}

func (cfg *DevConfig) FindPath(pathname string) (int, *DevPathConfig) {
	for index, path := range cfg.Paths {
		if path.FullPath == pathname {
			return index, &path
		}
	}
	return -1, nil
}

func (cfg *DevConfig) FindDevPath(pathname string) int {
	for index, path := range cfg.DevFolders {
		if path == pathname {
			return index
		}
	}
	return -1
}
func (pathConfig *DevPathConfig) Match(words []string) bool {
	lastIndex := -1
	for _, s := range words {
		if len(s) == 0 {
			continue
		}
		i := strings.Index(pathConfig.FullPath, s)
		if i <= lastIndex {
			return false
		}
		lastIndex = i
	}
	return true
}
