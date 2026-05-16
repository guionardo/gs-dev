package devservice

import (
	"log/slog"
	"os"
	"path"

	"github.com/guionardo/gs-dev/app/build"
	"github.com/guionardo/gs-dev/internal/consts"
	"github.com/guionardo/gs-dev/internal/fs_tools"
	projectdetector "github.com/guionardo/gs-dev/pkg/project_detector"
	"go.yaml.in/yaml/v3"
)

type (
	// LocalConfig is a configuration for a local folder
	LocalConfig struct {
		directory        string
		hasConfigFile    bool
		Ignore           bool   `yaml:"ignore"`
		IgnoreSubfolders bool   `yaml:"ignore_subfolders"`
		Description      string `yaml:"description"`
		ProjectType      string `yaml:"project_type"`
		ProjectName      string `yaml:"project_name"`
	}
)

const defaultConfigFile = "." + build.AppName

// NewLocalConfig creates a new local config from the directory
// If the directory does not exist, it returns an error
// If the directory contains a .gsdev.yaml file, it parses it and returns the local config
// If the directory does not contain a .gsdev.yaml file, it detects the project type, name and description from the directory
func NewLocalConfig(directory string) (*LocalConfig, error) {
	lc := &LocalConfig{
		directory: directory,
	}

	return lc, lc.Parse()
}

// Parse tries to parse the local config file or detect configuration from the directory
func (c *LocalConfig) Parse() error {
	filename, err := fs_tools.AssertFilename(path.Join(c.directory, defaultConfigFile))
	if err == nil {
		// Try to parse the local config file
		file, err2 := os.Open(filename) // nolint:gosec // G304 -- filename is validated
		if err2 == nil {
			defer file.Close() // nolint:errcheck // G307 -- defer error is ignored

			if err2 = yaml.NewDecoder(file).Decode(c); err2 != nil {
				return err2
			}

			c.hasConfigFile = true
		}
	}

	// Update project informations from the detected projects if not set in the config file
	if c.ProjectType == "" || c.ProjectName == "" || c.Description == "" {
		projectData, err := projectdetector.DetectProject(c.directory)
		if err == nil {
			slog.Debug("Project detected",
				slog.String("project", projectData.Name),
				slog.String("type", projectData.Type),
				slog.String("description", projectData.Description))
			c.ProjectName = projectData.Name
			c.ProjectType = projectData.Type
			c.Description = projectData.Description
		}
	}

	if !c.hasConfigFile {
		return nil
	}
	// Update description from the config file
	content, err := yaml.Marshal(c)
	if err == nil {
		err = os.WriteFile(filename, content, consts.FilesPermissions)
	}

	if err != nil {
		slog.Error("Error updating description from config file",
			slog.String("directory", c.directory),
			slog.Any("error", err))
	}

	return nil
}

func (c *LocalConfig) HasConfigFile() bool {
	return c.hasConfigFile
}

func (c *LocalConfig) Equal(other LocalConfig) bool {
	return c.Ignore == other.Ignore && c.IgnoreSubfolders == other.IgnoreSubfolders && c.Description == other.Description && c.ProjectType == other.ProjectType &&
		c.ProjectName == other.ProjectName
}
