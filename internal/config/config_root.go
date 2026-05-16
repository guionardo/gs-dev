package config

import (
	"errors"
	"log/slog"
	"os"
	"path"

	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/fs_tools"
	"github.com/guionardo/gs-dev/internal/logging"

	errs "github.com/guionardo/gs-dev/internal/errors"
	"gopkg.in/yaml.v3"
)

type (
	// ConfigRoot is an abstraction to group configurations by type
	ConfigRoot struct {
		configDir string
		data      map[string]any
		hash      map[string]uint64
	}
	Defaulter interface {
		Defaults()
	}
	ConfigType interface {
		Key() string
		Validate() error
	}
)

const defaultConfigFileExtension = ".yml"

func NewConfigRoot(configDir string) (config *ConfigRoot, err error) {
	configDir, err = fs_tools.AssertDirectory(configDir)
	if err != nil {
		logging.Debug("NewConfigRoot", slog.String("configDir", configDir), slog.Any("error", err))

		return nil, err
	}

	config = &ConfigRoot{
		configDir: configDir,
		data:      make(map[string]any),
		hash:      make(map[string]uint64),
	}

	logging.Debug("NewConfigRoot", slog.String("configDir", configDir))

	return config, nil
}

func (c ConfigRoot) Save() (err error) {
	// Check if the content of the config is the same as when loaded
	for key, data := range c.data {
		if loadedHash, ok := c.hash[key]; ok {
			actualHash := getHash(data.(ConfigType))
			if actualHash == loadedHash {
				continue
			}
		}

		configFileName := path.Join(c.configDir, key+defaultConfigFileExtension)
		if content, marshalErr := yaml.Marshal(data); marshalErr != nil {
			err = errors.Join(err, errs.NewError(marshalErr, "Error marshaling config %s - %w", false, key, marshalErr))
			continue
		} else if writeErr := os.WriteFile(configFileName, content, consts.FilesPermissions); writeErr != nil {
			err = errors.Join(err, errs.NewError(writeErr, "Error writing config %s - %w", false, writeErr))
			continue
		}
	}

	return err
}

// GetValue reads configuration from root. If the error is not null, the configuration is not recoverable
func GetValue[T ConfigType](config *ConfigRoot) (value T, err error) {
	anyValue, ok := config.data[value.Key()]
	if ok {
		if v, ok := anyValue.(T); ok {
			return v, nil
		}
	}

	configFileName := path.Join(config.configDir, value.Key()+defaultConfigFileExtension)

	value, err = readConfigFile[T](configFileName)
	if err == nil {
		// Success
		config.data[value.Key()] = value
		config.hash[value.Key()] = getHash(value)

		return value, nil
	}

	if e, ok := errors.AsType[errs.Error](err); ok {
		if e.IsRecoverable() {
			config.data[value.Key()] = value
			config.hash[value.Key()] = 0 // Force write default configuration

			return value, nil
		} else {
			logging.Warn("Configuration", slog.Any("error", e))
			return value, e
		}
	}

	logging.Error("Configuration", slog.Any("error", err))

	return value, err
}

func readConfigFile[T ConfigType](fileName string) (value T, err error) {
	content, err := os.ReadFile(fileName)
	if err == nil {
		err = yaml.Unmarshal(content, &value)
	}

	if err == nil {
		// configuration read with succes
		err = value.Validate()
	}

	if err == nil {
		return value, nil
	}

	// Failed on get data, try to set defaults
	var v any = &value
	if defaulter, ok := v.(Defaulter); ok {
		defaulter.Defaults()
		return value, errs.NewError(err, "Failed to read/parse config [%s]. Getting defaults. - %w", true, value.Key(), err)
	}

	return value, errs.NewConfigurationError("Failed to read/parse config [%s] - %w", value.Key(), err)
}

func SetValue[T ConfigType](config *ConfigRoot, value T) {
	config.data[value.Key()] = value
}
