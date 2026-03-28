package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"

	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/fs_tools"

	errs "github.com/guionardo/gs-dev/internal/errors"
	"gopkg.in/yaml.v3"
)

type (
	ConfigFile struct {
		fileName string
		data     map[string]any
	}
	Defaulter interface {
		Defaults()
	}
)

func NewConfigFile(filename string) (config *ConfigFile, err error) {
	filename, err = fs_tools.AssertFilename(filename)
	if err != nil {
		return nil, err
	}

	config = &ConfigFile{
		fileName: filename,
		data:     make(map[string]any),
	}

	content, err := os.ReadFile(filename)
	if errors.Is(err, fs.ErrNotExist) {
		return config, errs.NewError(err, "file does not exist", true)
	}

	if err != nil {
		return config, errs.NewError(err, "error reading file", false)
	}

	err = yaml.Unmarshal(content, &config.data)
	if err != nil {
		return config, errs.NewError(err, "error unmarshalling file", false)
	}

	return config, err
}

func (c ConfigFile) Save() (err error) {
	content, err := yaml.Marshal(c.data)
	if err == nil {
		err = os.WriteFile(c.fileName, content, consts.FilesPermissions)
	}

	return err
}

func (c *ConfigFile) SetValue(key string, value any) {
	c.data[key] = value
}

func GetValue[T any](config *ConfigFile, key string) (value T, err error) {
	v, ok := config.data[key]
	if !ok {
		slog.Debug("Key not found", slog.String("key", key))
		return value, errs.NewError(fmt.Errorf("key %s not found", key), "key not found", false)
	}

	content, err := yaml.Marshal(v)
	if err != nil {
		slog.Warn("Error marshalling value", slog.Any("error", err), slog.Any("value", v), slog.String("type", fmt.Sprintf("%T", v)))
		return value, fmt.Errorf("error marshalling value: %w", err)
	}

	err = yaml.Unmarshal(content, &value)
	if err != nil {
		slog.Warn("Error unmarshalling value", slog.Any("error", err), slog.Any("value", content), slog.String("type", fmt.Sprintf("%T", value)))
		return value, errs.NewError(fmt.Errorf("error unmarshalling value: %w", err), "error unmarshalling value", false)
	}

	slog.Debug("Value unmarshalled", slog.Any("value", value), slog.String("type", fmt.Sprintf("%T", value)))

	var vAny any = value
	if v, ok := vAny.(Defaulter); ok {
		v.Defaults()
	}

	return value, nil
}

func SetValue[T any](config ConfigFile, key string, value T) (err error) {
	config.data[key] = value
	return nil
}
