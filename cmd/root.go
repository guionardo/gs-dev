/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"os"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/console"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     app.ToolName,
	Short:   app.ShortDescription,
	Long:    app.Description,
	Version: app.Version,
	PreRun: func(cmd *cobra.Command, args []string) {
		console.OutputNeutral("%s v%s\n", app.ToolName, app.Version)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {	
	err := rootCmd.ExecuteContext(configs.GetContextConfiguration())
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	configs.SetupConfigurationRoot(rootCmd)
}
