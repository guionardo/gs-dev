package cmd

import (
	"github.com/guionardo/gs-dev/app/setup"
	"github.com/guionardo/gs-dev/configs"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Application configuration",
	Run: func(cmd *cobra.Command, args []string) {
		setup.Setup(cmd)
		cfg := configs.ValidateConfiguration(cmd)
		cmd.Printf("setup called %v", cfg)
	},
}

func init() {
	rootCmd.AddCommand(setupCmd)
}
