package configurations

import (
	"github.com/google/uuid"
	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/internal/errors"
)

type (
	PadServerConfig struct {
		Port             int    `yaml:"port" default:"8080"`
		APIKey           string `yaml:"api_key" default:""`
		StorageDirectory string `yaml:"storage_directory" default:""`
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

	if len(c.StorageDirectory) == 0 {
		return errors.NewError(nil, "storage directory cannot be empty", false)
	}

	if !pathtools.DirExists(c.StorageDirectory) {
		return errors.NewError(nil, "storage directory does not exist: %s", false, c.StorageDirectory)
	}

	return nil
}

func (c PadServerConfig) Key() string {
	return "pad_server"
}
