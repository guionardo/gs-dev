package padservice

import (
	"github.com/google/uuid"
	"github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/pad/storage"
)

type (
	PadServerConfig struct {
		Port          int                             `yaml:"port" default:"8080"`
		APIKey        string                          `yaml:"api_key" default:""`
		StorageConfig storage.FileSystemStorageConfig `yaml:"storage"`
	}
)

func NewPadServerConfig() *PadServerConfig {
	config := &PadServerConfig{}
	config.Defaults()

	return config
}

func (c *PadServerConfig) Defaults() {
	if c != nil && (c.Port <= 0 || c.Port > 65535) {
		c.Port = 8080
	}

	if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
		c.APIKey = ""
	}

	c.StorageConfig.Defaults()
}

func (c PadServerConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return errors.NewError(nil, "invalid port number: %d", true, c.Port)
	}

	if c.APIKey != "" {
		if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
			return errors.NewError(nil, "invalid API key: %s", true, c.APIKey)
		}
	}

	return c.StorageConfig.Validate()
}

func (c PadServerConfig) Key() string {
	return "pad_server"
}
