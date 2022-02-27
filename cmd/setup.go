/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"github.com/guionardo/gs-dev/configs"
	"github.com/guionardo/gs-dev/internal/console"
	"github.com/spf13/cobra"
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Application configuration",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := configs.ValidateConfiguration(cmd)
		if cfg.ErrorCode > 0 {
			console.OutputError("Missing setup: %s", cfg.Error)
		}
		cmd.Help()

	},
}

func init() {
	rootCmd.AddCommand(setupCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// setupCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// setupCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
