package config

import "path"

type PadConfig struct {
	fileName   string
	BackendUrl string `yaml:"backend_url"`
	ApiKey     string `yaml:"api_key"`
}

func NewPadConfig(root string) *PadConfig {
	return &PadConfig{
		fileName: path.Join(root, "pad.yaml"),
	}
}
