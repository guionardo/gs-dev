package configs

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/guionardo/gs-dev/app"
	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type RootConfig struct {
	DataFolder string    `yaml:"-"`
	ConfigFile string    `yaml:"-"`
	ErrorCode  int       `yaml:"-"`
	Error      string    `yaml:"-"`
	DevConfig  DevConfig `yaml:"dev_config"`
}

const (
	GSDEV_CONFIGURATION_ENV   = "GSDEV_CONFIGURATION"
	NO_ERROR                  = 0
	ERROR_MISSING_CONF_PATH   = 1
	ERROR_CONF_PATH_NOT_FOUND = 2
	ERROR_CONF_FILE_NOT_FOUND = 3
	ERROR_CONF_FILE_READ      = 4
	ERROR_CONF_FILE_INVALID   = 5
)

var DefaultConfigurationPath string
var ConfigurationPath string

func init() {
	homeDir, _ := os.UserHomeDir()
	ConfigurationPath = os.Getenv(GSDEV_CONFIGURATION_ENV)
	if len(ConfigurationPath) == 0 {
		DefaultConfigurationPath = path.Join(homeDir, fmt.Sprintf(".%s", app.ToolName))
	} else {
		DefaultConfigurationPath = ConfigurationPath
	}
}

func SetupConfigurationRoot(rootCommand *cobra.Command) {
	rootCommand.PersistentFlags().StringVarP(&ConfigurationPath, "config", "c",
		DefaultConfigurationPath, fmt.Sprintf("configuration path (ENV=%s)", GSDEV_CONFIGURATION_ENV))
}

func (cfg *RootConfig) SetErrorf(errorCode int, format string, a ...interface{}) {
	cfg.ErrorCode = errorCode
	if errorCode == 0 {
		cfg.Error = ""
	} else {
		cfg.Error = fmt.Sprintf("#%d %s", errorCode, fmt.Sprintf(format, a...))
	}
}

func ValidateConfiguration(cmd *cobra.Command) (cfg RootConfig) {
	cfg = RootConfig{}
	dataFolder, err := cmd.Flags().GetString("config")
	cfg.DataFolder = dataFolder
	if err != nil {
		cfg.SetErrorf(ERROR_MISSING_CONF_PATH, "MISSING CONFIGURATION PATH")
		return
	}
	if _, err = os.Stat(dataFolder); err != nil {
		cfg.SetErrorf(ERROR_CONF_PATH_NOT_FOUND, "CONFIGURATION PATH NOT FOUND %s", dataFolder)
		return
	}
	configFile := path.Join(dataFolder, fmt.Sprintf("%s.yaml", app.ToolName))
	cfg.ConfigFile = configFile
	if _, err = os.Stat(configFile); err != nil {
		cfg.SetErrorf(ERROR_CONF_FILE_NOT_FOUND, "CONFIGURATION FILE NOT FOUND %s", configFile)
		return
	}
	err = cfg.ReadFile(configFile)
	if err != nil {
		cfg.SetErrorf(ERROR_CONF_FILE_READ, "CONFIGURATION FILE FAILED TO READ %s - %v", configFile, err)
		return
	}
	return
}

func (cfg *RootConfig) ReadFile(filename string) error {
	file, err := os.ReadFile(filename)
	if err == nil {
		err = yaml.Unmarshal(file, &cfg)
	}
	return err
}

func (cfg *RootConfig) WriteFile() error {
	bytes, err := yaml.Marshal(cfg)
	if err == nil {
		err = os.WriteFile(cfg.ConfigFile, bytes, 0666)
	}
	return err
}

func (cfg *RootConfig) Prompt() bool {
	if !Confirm("Setup basic configuration?", true) {
		return false
	}
	dataFolder := cfg.DataFolder
	if len(dataFolder) == 0 {
		dataFolder = DefaultConfigurationPath
	}
	dataFolder, err := PromptPath("Data folder", dataFolder)
	if err == nil {
		cfg.DataFolder = dataFolder
		return true
	}
	return false
}

func ReadRootConfiguration() (cfg *RootConfig) {
	cfg = &RootConfig{}

	cfg.DataFolder = os.Getenv(GSDEV_CONFIGURATION_ENV)
	if len(cfg.DataFolder) == 0 {
		homeDir, _ := os.UserHomeDir()
		cfg.DataFolder = path.Join(homeDir, fmt.Sprintf(".%s", app.ToolName))
	}

	if stat, err := os.Stat(cfg.DataFolder); err != nil || !stat.IsDir() {
		cfg.SetErrorf(ERROR_CONF_PATH_NOT_FOUND, "PATH NOT FOUND %s", cfg.DataFolder)
	}
	cfg.ConfigFile = path.Join(cfg.DataFolder, fmt.Sprintf("%s.yaml", app.ToolName))
	if stat, err := os.Stat(cfg.ConfigFile); err != nil || stat.IsDir() {
		cfg.SetErrorf(ERROR_CONF_FILE_NOT_FOUND, "CONFIG FILE NOT FOUND %s", cfg.ConfigFile)
	} else {
		if err := cfg.ReadFile(cfg.ConfigFile); err != nil {
			cfg.SetErrorf(ERROR_CONF_FILE_INVALID, "ERROR READING CONFIG FILE %s - %s", cfg.ConfigFile, err)
		}
	}
	return
}

type cfgKey string

const CfgKey = cfgKey("cfg")

func GetContextConfiguration() context.Context {
	cfg := ReadRootConfiguration()
	return context.WithValue(context.Background(), CfgKey, cfg)
}

func GetConfiguration(cmd *cobra.Command) (cfg *RootConfig, err error) {
	contextCfg := cmd.Context().Value(CfgKey)

	if contextCfg == nil {
		cfg = nil
		err = fmt.Errorf("no root configuration")
	} else {
		cfg = contextCfg.(*RootConfig)
		if cfg.ErrorCode > 0 {
			err = fmt.Errorf(cfg.Error)
		}
	}
	return
}

func (cfg *RootConfig) ValidateDataFolder() error {
	if stat, err := os.Stat(cfg.DataFolder); err == nil && stat.IsDir() {
		return nil
	}
	return pathtools.CreatePath(cfg.DataFolder)
}

func (cfg *RootConfig) ValidateConfigFile() error {
	if stat, err := os.Stat(cfg.ConfigFile); err == nil && !stat.IsDir() {
		return nil
	}
	return cfg.WriteFile()
}
