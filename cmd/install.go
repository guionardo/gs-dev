package cmd

import (
	"github.com/guionardo/gs-dev/app/install"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install bindings on your shell profile",
	RunE:  installRun,
}

func init() {
	installCmd.Flags().BoolP("uninstall", "u", false, "Uninstall bindings")
	rootCmd.AddCommand(installCmd)
}

func installRun(cmd *cobra.Command, args []string) error {
	if uninstall, err := cmd.Flags().GetBool("uninstall"); err != nil {
		return err
	} else {
		if uninstall {
			return install.RunUninstall()
		}
		return install.RunInstall()
	}
}
