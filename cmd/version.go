package cmd

import (
	"github.com/guionardo/gs-dev/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Application version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("%s v%s\n", app.ToolName, app.Version)
		},
	})
}
