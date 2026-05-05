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
	// ConfigRoot is an abstraction to group configurations by type
	ConfigRoot struct {
		fileName string
		data     map[string]any
	}
	Defaulter interface {
		Defaults()
	}
	ConfigType interface {
		Key() string
	}
	Validater interface {
		Validate() error
	}
)

func NewConfigFile(filename string) (config *ConfigRoot, err error) {
	filename, err = fs_tools.AssertFilename(filename)
	if err != nil {
		return nil, err
	}

	config = &ConfigRoot{
		fileName: filename,
		data:     make(map[string]any),
	}

	content, err := os.ReadFile(filename) // nolint: gosec // file validated in the AssertFileName function
	if errors.Is(err, fs.ErrNotExist) {
		return config, errs.NewError(err, "file does not exist", true)
	}

	if err != nil {
		return config, errs.NewError(err, "error reading configuration file: %s", false, filename)
	}

	err = yaml.Unmarshal(content, &config.data)
	if err != nil {
		return config, errs.NewError(err, "error unmarshalling file", false)
	}

	return config, err
}

func (c ConfigRoot) Save() (err error) {
	content, err := yaml.Marshal(c.data)
	if err == nil {
		err = os.WriteFile(c.fileName, content, consts.FilesPermissions)
	}

	return err
}

func GetValue[T ConfigType](config *ConfigRoot) (value T, err error) {
	key := value.Key()

	v, ok := config.data[key]
	if !ok {
		slog.Debug("Required configuration not found", slog.String("key", key))

		return value, errs.NewError(nil, "key %s not found", false, key)
	}

	content, _ := yaml.Marshal(v)

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

	if v, ok := vAny.(Validater); ok {
		err = v.Validate()
	}

	return value, err
}

func SetValue[T ConfigType](config *ConfigRoot, value T) {
	config.data[value.Key()] = value
}
