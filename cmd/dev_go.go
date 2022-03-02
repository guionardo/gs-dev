/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"github.com/guionardo/gs-dev/app/dev"
	"github.com/spf13/cobra"
)

var goCmd = &cobra.Command{
	Use:   "go",
	Short: "Projects paths navigator",
	Args:  cobra.MinimumNArgs(1),
	Long:  `add words for search for path`,
	Run: func(cmd *cobra.Command, args []string) {
		dev.DevGo(cmd, args)
	},
}

func init() {
	goCmd.Flags().StringP("output", "o", "", "Output shell file for command execution")
	devCmd.AddCommand(goCmd)
}
