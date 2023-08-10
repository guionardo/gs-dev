// CODE GENERATED. DO NOT EDIT.
package config

import (
	"os"	

	"gopkg.in/yaml.v3"
)

func (cfg *DevConfig) Save() error {
	if content, err := yaml.Marshal(cfg); err != nil {
		return err
	} else {
		return os.WriteFile(cfg.fileName, content, 0644)
	}
}

func (cfg *DevConfig) Load() error {
	if content, err := os.ReadFile(cfg.fileName); err != nil {
		return err
	} else {
		return yaml.Unmarshal(content, cfg)
	}
}
