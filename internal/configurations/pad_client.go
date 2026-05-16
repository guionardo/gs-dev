package configurations

import (
	"net/url"

	"github.com/google/uuid"
	"github.com/guionardo/gs-dev/internal/errors"
)

type PadClientConfig struct {
	Enabled    bool   `yaml:"enabled" default:"true"`
	BackendURL string `yaml:"backend_url" default:"http://localhost:8080"`
	APIKey     string `yaml:"api_key" default:""`
}

func NewPadClientConfig() *PadClientConfig {
	config := &PadClientConfig{}
	config.Defaults()

	return config
}

func (c *PadClientConfig) Defaults() {
	if beURL, err := url.Parse(c.BackendURL); err != nil || beURL.Scheme != "http" && beURL.Scheme != "https" {
		c.BackendURL = ""
		c.Enabled = false
	}

	if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
		c.APIKey = ""
		c.Enabled = false
	}
}

func (c PadClientConfig) Validate() error {
	if !c.Enabled {
		return nil
	}

	backendURL, err := url.Parse(c.BackendURL)
	if err != nil || backendURL.Scheme != "http" && backendURL.Scheme != "https" {
		return errors.NewError(nil, "invalid backend URL: %s", false, c.BackendURL)
	}

	// API Key can be empty, but if it's not, it should be a valid UUID
	if c.APIKey != "" {
		if akUUID, err := uuid.Parse(c.APIKey); err != nil || akUUID == uuid.Nil {
			return errors.NewError(nil, "invalid API key: %s", false, c.APIKey)
		}
	}

	return nil
}

func (c PadClientConfig) Key() string {
	return "pad_client"
}
