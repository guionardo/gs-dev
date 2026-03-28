package dev

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	pathtools "github.com/guionardo/go/path_tools"
	"github.com/guionardo/gs-dev/pkg/tools/files"
	"gopkg.in/yaml.v3"
)

type LocalConfig struct {
	directory        string
	Ignore           bool   `yaml:"ignore"`
	IgnoreSubfolders bool   `yaml:"ignore_subfolders"`
	Description      string `yaml:"description"`
}

func NewLocalConfig(directory string) (*LocalConfig, error) {
	lc := &LocalConfig{
		directory: filepath.Clean(directory),
	}

	return lc, lc.Parse()
}

func (c *LocalConfig) Parse() error {
	if !pathtools.DirExists(c.directory) {
		return os.ErrNotExist
	}

	// Find the .gsdev directory in the current directory
	if c.parseDirectory() == nil {
		return nil
	}
	// Try to parse the gsdev file in the directory
	return c.parseGsDevFile(c.directory)
}

func (c *LocalConfig) parseDirectory() error {
	dir := path.Join(c.directory, ".gsdev")
	if !pathtools.DirExists(dir) {
		return os.ErrNotExist
	}

	return c.parseGsDevFile(dir)
}

func (c *LocalConfig) parseGsDevFile(dir string) error {
	gsDevFile := files.FindFirst(dir, ".gsdev.yaml", ".gsdev.yml", ".gs_dev.yaml", ".gs_dev.yml", ".gs-dev.yaml", ".gs-dev.yml", ".gsdev", ".gs-dev", ".gs_dev")
	if gsDevFile == "" {
		return fmt.Errorf("gs-dev local config file not found in %s", dir)
	}

	content, err := os.ReadFile(filepath.Clean(gsDevFile))
	if err != nil {
		return fmt.Errorf("error reading gs-dev local config file: %s - %w", gsDevFile, err)
	}

	if err = yaml.Unmarshal(content, c); err != nil {
		return fmt.Errorf("error unmarshalling gs-dev local config file: %s - %w", gsDevFile, err)
	}

	return nil
}
