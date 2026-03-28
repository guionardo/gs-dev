package commons

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/pkg/plugins"

	"gopkg.in/yaml.v3"
)

type PluginConfiguration struct {
	Enabled bool `yaml:"enabled"`
}

var configFilenamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func LoadConfiguration(plugin plugins.CliPlugin, configuration any, required bool, filename ...string) error {
	configFile, err := resolveConfigurationPath(plugin, filename...)
	if err != nil {
		return err
	}

	content, err := os.ReadFile(configFile) // #nosec G304 -- constrained basename under config directory.
	if os.IsNotExist(err) && !required {
		return nil
	}

	if err == nil {
		err = yaml.Unmarshal(content, configuration)
	}

	return err
}

func SaveConfiguration(plugin plugins.CliPlugin, configuration any, filename ...string) error {
	configFile, err := resolveConfigurationPath(plugin, filename...)
	if err != nil {
		return err
	}

	content, err := yaml.Marshal(configuration)
	if err == nil {
		err = os.WriteFile(configFile, content, consts.FilesPermissions) // #nosec G304 -- constrained basename under config directory.
	}

	return err
}

func resolveConfigurationPath(plugin plugins.CliPlugin, filename ...string) (string, error) {
	configurationName := plugin.GetConfiguration().Name
	if len(filename) == 0 {
		filename = []string{configurationName}
	}

	file := filename[0]
	if !configFilenamePattern.MatchString(file) {
		return "", fmt.Errorf("invalid configuration filename %q", file)
	}

	return filepath.Join(config.GetConfigDir(), file+".yaml"), nil
}
