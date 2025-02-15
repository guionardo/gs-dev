package commons

import (
	"os"
	"path"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/pkg/plugins"

	"gopkg.in/yaml.v3"
)

type PluginConfiguration struct {
	Enabled bool `yaml:"enabled"`
}

func LoadConfiguration(plugin plugins.CliPlugin, configuration interface{}, required bool, filename ...string) error {
	if len(filename) == 0 {
		filename = []string{plugin.GetConfiguration().Name}
	}
	configFile := path.Join(config.GetConfigDir(), filename[0]+".yaml")
	content, err := os.ReadFile(configFile)
	if os.IsNotExist(err) && !required {
		return nil
	} else if err == nil {
		err = yaml.Unmarshal(content, configuration)
	}
	return err
}

func SaveConfiguration(plugin plugins.CliPlugin, configuration interface{}, filename ...string) error {
	if len(filename) == 0 {
		filename = []string{plugin.GetConfiguration().Name}
	}
	configFile := path.Join(config.GetConfigDir(), filename[0]+".yaml")
	content, err := yaml.Marshal(configuration)
	if err == nil {
		err = os.WriteFile(configFile, content, 0644)
	}
	return err
}
