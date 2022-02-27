package setup

import (
	"os"
	"path"
	"time"

	"github.com/guionardo/gs-dev/app/dev"
	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/console"
	"github.com/spf13/cobra"
)

func SetupInit(cmd *cobra.Command) error {
	cfg, err := configs.GetConfiguration(cmd)

	if err == nil {
		if force, err := cmd.Flags().GetBool("force"); err != nil || !force {
			console.OutputSuccess("Setup is OK - %s\n", cfg.DataFolder)
			return nil
		}
	}
	console.OutputInfo("Validating data folder: %s\n", cfg.DataFolder)
	if err = cfg.ValidateDataFolder(); err != nil {
		return err
	}
	console.OutputInfo("Validating configuration file: %s\n", cfg.ConfigFile)
	if err = cfg.ValidateConfigFile(); err != nil {
		return err
	}

	home, _ := os.UserHomeDir()
	devFolder := path.Join(home, "dev")

	console.OutputInfo("Validating dev folders")

	if stat, err := os.Stat(devFolder); err == nil {
		if stat.IsDir() {
			if index := cfg.DevConfig.FindDevPath(devFolder); index == -1 {
				cfg.DevConfig.DevFolders = append(cfg.DevConfig.DevFolders, devFolder)
				console.OutputNeutral("+ Added default: %s\n", devFolder)
			}
		}
	}
	maxSubLevels, err := configs.PromptInt("Maximum sub levels", 3, 1, 10)
	if err != nil {
		maxSubLevels = 3
		// return err
	}
	cfg.DevConfig.MaxSubLevels = maxSubLevels

	updateInterval, err := configs.PromptInt("Update interval (seconds)", 86400, 60, 604800)
	if err != nil {
		updateInterval = 86400
	}
	cfg.DevConfig.UpdateInterval = time.Duration(updateInterval * int(time.Second))
	console.OutputNeutral("Update interval = %v\n", cfg.DevConfig.UpdateInterval)
	changes, err := dev.UpdatePaths(cfg)
	if err != nil {
		console.OutputError("Failed to update paths: %v", err)
	} else {
		console.OutputSuccess("Changed paths: %d\n", changes)
	}

	if err := cfg.WriteFile(); err != nil {
		console.OutputError("  ! FAILED: %v\n", err)
		return err
	}
	console.OutputSuccess(" Success")

	return nil

}

func promptCfg(cfg *configs.RootConfig, cmd *cobra.Command) {
	if cfg.DevConfig.Prompt() {
		return
	}

	configurationPath, err := console.PromptFolder("Configuration path", configs.DefaultConfigurationPath)

	if err != nil {
		cmd.PrintErrln(err)
	}
	cfg.DataFolder = configurationPath
}
