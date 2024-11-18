package configs

import (
	"fmt"
	"os"
	"path"

	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type (
	ConfigError struct {
		ErrorCode int   `yaml:"-"`
		Error     error `yaml:"-"`
	}
	ConfigFolder struct {
		ConfigError `yaml:"-"`
		DataFolder  string `yaml:"-"`
	}
)

func (cfg *ConfigError) SetErrorf(errorCode int, format string, a ...interface{}) {
	cfg.ErrorCode = errorCode
	if errorCode == 0 {
		cfg.Error = nil
	} else {
		cfg.Error = fmt.Errorf("#%d %s", errorCode, fmt.Sprintf(format, a...))
	}
}

func (cfg *ConfigFolder) LoadFile(fileName string) (content []byte, err error) {
	return os.ReadFile(path.Join(cfg.DataFolder, fileName+".yaml"))
}

func (cfg *ConfigFolder) SaveFile(fileName string, data any) error {
	if content, err := yaml.Marshal(data); err == nil {
		return os.WriteFile(path.Join(cfg.DataFolder, fileName+".yaml"), content, 0644)
	} else {
		return err
	}
}

func NewConfigFolder(cmd *cobra.Command) (cfg *ConfigFolder) {
	cfg = &ConfigFolder{}

	if dataFolder, err := cmd.Flags().GetString("config"); err == nil {
		cfg.DataFolder = dataFolder
	} else {
		cfg.SetErrorf(ERROR_MISSING_CONF_PATH, "MISSING CONFIGURATION PATH")
		return
	}

	if stat, err := os.Stat(cfg.DataFolder); err != nil {
		if os.IsNotExist(err) {
			if err = pathtools.CreatePath(cfg.DataFolder); err != nil {
				cfg.SetErrorf(ERROR_CONF_PATH_IO, "CONFIGURATION PATH CANNOT BE CREAED %v", err)
				return
			}
			cmd.Printf("Configuration path created: %s", cfg.DataFolder)
		}
	} else if !stat.IsDir() {
		cfg.SetErrorf(ERROR_CONF_PATH_IO, "CONFIGURATION PATH IS A FILE: %s", cfg.DataFolder)
		return
	}

	return
}
