/*
Copyright © 2022 Guionardo Furlan <guionardo@gmail.com>

*/
package cmd

import (
	"github.com/guionardo/gs-dev/app/notify"
	"github.com/spf13/cobra"
)

// notifyAddCmd represents the notifyAdd command
var notifyAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new message notification",
	Args:  cobra.MinimumNArgs(1),
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return notify.AddNotify(cmd,args)

	},
}

func init() {		
	notifyAddCmd.Flags().StringP("type", "t", "info", "Message type [info | warning | error]")
	notifyCmd.AddCommand(notifyAddCmd)
}
