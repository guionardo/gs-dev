package setup

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/configs"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func SetupInit(cmd *cobra.Command) {
	cmd.Println(app.ToolName)
	cfg := configs.ValidateConfiguration(cmd)
	if cfg.ErrorCode == configs.ERROR_MISSING_CONF_PATH {
		cmd.PrintErrf("%s", cfg.Error)
		return
	}
	if !cfg.Prompt(){
		return
	}
	if cfg.ErrorCode == configs.NO_ERROR {
		cmd.Printf("%s Setup is ok on %s", promptui.IconGood, cfg.DataFolder)
		return
	}
	if cfg.ErrorCode == configs.ERROR_CONF_PATH_NOT_FOUND {
		cmd.Printf("+ Creating configuration path: %s\n", configs.ConfigurationPath)
		syscall.Umask(0)
		if err := os.Mkdir(configs.ConfigurationPath, os.ModeSticky|os.ModePerm); err != nil {
			cmd.PrintErrf("  ! FAILED: %v\n", err)
			return
		}
		if err := os.Chmod(configs.ConfigurationPath, os.ModeSticky|os.ModePerm); err != nil {
			cmd.PrintErrf("  ! FAILED: %v\n", err)
			return
		}
		cfg.ErrorCode = configs.ERROR_CONF_FILE_NOT_FOUND
		cfg.ConfigFile = path.Join(cfg.DataFolder, fmt.Sprintf("%s.yaml", app.ToolName))
	}

	if cfg.ErrorCode == configs.ERROR_CONF_FILE_NOT_FOUND {
		cmd.Printf("+ Creating configuration file: %s\n", cfg.ConfigFile)
		promptCfg(&cfg, cmd)

		if err := cfg.WriteFile(cfg.ConfigFile); err != nil {
			cmd.PrintErrf("  ! FAILED: %v\n", err)
			return
		}
	}

	cmd.Println(promptui.IconGood + " Success")
}

func promptCfg(cfg *configs.RootConfig, cmd *cobra.Command) {
	if cfg.DevConfig.Prompt(){
		return
	}
	validateFolder := func(folderName string) error {
		if _, err := os.Stat(folderName); err != nil {
			return errors.New("path not found")
		}
		return nil
	}
	prompt := promptui.Prompt{Label: "Configuration path",
		Validate: validateFolder,
		Default:  configs.DefaultConfigurationPath,
	}
	configurationPath, err := prompt.Run()
	if err != nil {
		cmd.PrintErrln(err)
	}
	cfg.DataFolder = configurationPath
}
