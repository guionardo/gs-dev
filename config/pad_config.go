package config

import (
	"path"
	"sync"
)

type PadConfig struct {
	fileName   string
	lock       sync.Mutex
	BackendUrl string `yaml:"backend_url"`
	ApiKey     string `yaml:"api_key"`
}

func NewPadConfig(root string) *PadConfig {
	return &PadConfig{
		fileName: path.Join(root, "pad.yaml"),
	}
}
