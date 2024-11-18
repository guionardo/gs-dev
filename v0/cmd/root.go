/*
Copyright © 2023 Guionardo Furlan <guionardo@gmail.com>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/app/setup"
	"github.com/guionardo/gs-dev/config"
	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/logger"
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
	cmd, _, err := rootCmd.Find(os.Args[1:])
	// default cmd if no cmd is given
	if err != nil && cmd.Use == rootCmd.Use {
		// && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp
		args := append([]string{devCmd.Use}, os.Args[1:]...)
		rootCmd.SetArgs(args)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// err := rootCmd.Execute()
	// if err != nil {
	// 	os.Exit(1)
	// }
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	configs.SetupConfigurationRoot(rootCmd)

	setup.SetRootCmd(rootCmd)
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	// Check if config root is ok
	if repositoryFolder, err := cmd.Flags().GetString("config"); err == nil {
		if err = config.ValidateRepositoryFolder(repositoryFolder); err != nil {
			return err
		}
	}

	// Debug
	if debug, err := cmd.Flags().GetBool("debug"); err == nil {
		logger.SetDebugMode(debug)
		logger.Debug("DEBUG MODE ON")
	}

	return nil
}
