package configs

import (
	"os"
	"path"
	"testing"

	pathtools "github.com/guionardo/gs-dev/internal/path_tools"
)

func createTree(basePath string) error {
	paths := []string{"alpha", "beta", "gama"}
	subpaths := []string{"mercur", "venus", "earth"}
	for _, p := range paths {
		for _, sp := range subpaths {
			x := path.Join(basePath, p, sp)
			if err := pathtools.CreatePath(x); err != nil {
				return err
			}
		}
	}
	return nil
}

func TestRootConfig_SaveToFile(t *testing.T) {
	tempDir := path.Join(os.TempDir(), "gs_dev_test")
	if err := pathtools.CreatePath(tempDir); err != nil {
		t.Errorf("Failed to create temporary path %v", tempDir)
	}
	defer os.RemoveAll(tempDir)

	devPath := path.Join(os.TempDir(), "dev")
	if err := pathtools.CreatePath(devPath); err != nil {
		t.Errorf("Failed to create temporary path %v", devPath)
	}
	defer os.RemoveAll(devPath)

	if err := createTree(devPath); err != nil {
		t.Errorf("Failed to create tree %v", err)
	}

	t.Run("Saving", func(t *testing.T) {
		config := &RootConfig{
			DataFolder: tempDir,
			ConfigFile: path.Join(tempDir, CONFIGURATION_FILE),
			ErrorCode:  0,
			Error:      "",
			DevConfig:  &DevConfig{},
		}
		config.DevConfig.DevFolders.
		if err := config.SaveToFile(config.ConfigFile); err != nil {
			t.Errorf("Failed to save configuration - %v", err)
		}
	})
	t.Run("Loading", func(t *testing.T) {
		config := LoadConfiguration(tempDir)
		if config.ErrorCode != 0 {
			t.Errorf("Failed to load configuration - %s", config.Error)
		}
	})

}
