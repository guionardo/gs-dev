package cmd

import (
	initshell "github.com/guionardo/gs-dev/app/init_shell"
	"github.com/spf13/cobra"
)

func init() {
	initCommand := &cobra.Command{
		Use:   "init",
		Short: "Initialization for shell alias",
		RunE: func(cmd *cobra.Command, args []string) error {
			return initshell.InitAction()
		},
	}
	rootCmd.AddCommand(initCommand)

}
