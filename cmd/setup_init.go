/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"github.com/guionardo/gs-dev/app/setup"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize configuration and data repository",
	RunE: func(cmd *cobra.Command, args []string) error {
		return setup.SetupInit(cmd)
	},
}

func init() {
	initCmd.Flags().BoolP("force", "f", false, "Force reinitialization")
	setupCmd.AddCommand(initCmd)
}
