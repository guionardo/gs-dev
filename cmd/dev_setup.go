/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"github.com/guionardo/gs-dev/app/dev"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var goSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup for the current path",
	RunE: func(cmd *cobra.Command, args []string) error {
		return dev.DevSetup(cmd)
	},
}

func init() {
	empty := make([]string, 0)
	goSetupCmd.Flags().StringP("status", "e", "", "Options: enabled | disabled")
	goSetupCmd.Flags().StringArrayP("after-command", "a", empty, "After command(s)")
	goSetupCmd.Flags().BoolP("ignore-subpaths", "i", false, "Ignore sub paths")
	devCmd.AddCommand(goSetupCmd)

}
