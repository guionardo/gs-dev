package dev

import (
	"reflect"
	"testing"

	"github.com/guionardo/gs-dev/configs"
	"github.com/spf13/cobra"
)

func TestDevConfig_Save_And_Load(t *testing.T) {

	t.Run("Default", func(t *testing.T) {
		cmd := &cobra.Command{}
		tmp := t.TempDir()
		configs.SetupConfigurationRoot(cmd)
		cmd.SetArgs([]string{"--config", tmp})
		cmd.Execute()

		cf := configs.NewConfigFolder(cmd)

		if cf.ConfigError.Error != nil {
			t.Errorf("NewConfigFolder error = %v", cf.ConfigError.Error)
			return
		}

		cfg := NewDevConfig(cf)
		cfg.DevFolders = []string{"/tmp", "/home", "/test"}
		cfg.MaxSubLevels = 3
		if err := cfg.Save(); err != nil {
			t.Errorf("DevConfig.Save() error = %v", err)
		}

		cfg2 := NewDevConfig(cf)
		if err := cfg2.Load(); err != nil {
			t.Errorf("DevConfig.Load() error = %v", err)
		}
		if !reflect.DeepEqual(cfg.DevFolders, cfg2.DevFolders) {
			t.Errorf("Expected dev folders %v -> got %v", cfg.DevFolders, cfg2.DevFolders)
		}
	})

}
