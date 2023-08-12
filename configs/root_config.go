// TODO: Verificar se root_config.go é necessário
package configs

import (
	"fmt"
	"os"
	"path"

	"github.com/guionardo/gs-dev/app"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type RootConfig struct {
	DataFolder string `yaml:"-"`
	ConfigFile string `yaml:"-"`
	ErrorCode  int    `yaml:"-"`
	Error      string `yaml:"-"`
}

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
	rootCommand.PersistentFlags().StringVar(&ConfigurationPath, "config",
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
		err = yaml.Unmarshal(file, cfg)
	}
	return err
}

func (cfg *RootConfig) WriteFile(filename string) error {
	bytes, err := yaml.Marshal(cfg)
	if err == nil {
		err = os.WriteFile(filename, bytes, 0666)
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
