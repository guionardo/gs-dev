/*
Copyright © 2023 Guionardo Furlan <guionardo@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/app/setup"
	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/configs"
	"github.com/spf13/cobra"
)

//go:generate go run ../gen/docs.go

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:               app.ToolName,
	Short:             app.ShortDescription,
	Long:              app.Description,
	PersistentPreRunE: rootPreRun,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	configs.SetupConfigurationRoot(rootCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	setup.SetRootCmd(rootCmd)
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	// Check if config root is ok
	if repositoryFolder, err := cmd.Flags().GetString("config"); err == nil {
		if err = config.ValidateRepositoryFolder(repositoryFolder); err != nil {
			return err
		}
	}

	return nil
}
