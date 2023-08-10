package cmd

import (
	"errors"
	"os"
	"path"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/app/dev"
	"github.com/spf13/cobra"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Rapid access to your development folders",
	Long: `This feature allows you to access your projects folder
quickly and get information about it`,
	RunE: devRun,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	addCommand := &cobra.Command{
		Use:   "add",
		Short: "Add root folder",
		Args:  cobra.ExactArgs(1),
		RunE:  devRunAdd,
	}
	addCommand.Flags().Uint8P("max-sublevels", "m", 3, "Maximum sub level for folders")

	removeComand := &cobra.Command{
		Use:   "remove",
		Short: "Remove root folder",
		Args:  cobra.ExactArgs(1),
		RunE:  devRunRemove,
	}

	syncCommand := &cobra.Command{
		Use:   "sync",
		Short: "Force synchronization of folders",
		RunE:  devRunSync,
	}

	ignoreChildCommand := &cobra.Command{
		Use:     "child",
		Short:   "Enable/disable children",
		Long:    "use arguments 'enable' or 'disable'",
		Args:    cobra.ExactArgs(1),
		RunE:    devRunChild,
		Example: "dev child --enable /home/some/folder",
	}
	ignoreChildCommand.Flags().BoolP("enable", "e", false, "Enable children")
	ignoreChildCommand.Flags().BoolP("disable", "d", false, "Disable children")
	ignoreChildCommand.MarkFlagsMutuallyExclusive("enable", "disable")

	devCmd.AddCommand(
		addCommand,
		removeComand,
		syncCommand,
		ignoreChildCommand,
	)

	output := path.Join(os.TempDir(), app.ToolName)
	devCmd.Flags().StringP("output", "o", output, "Output script for shell alias")
	rootCmd.AddCommand(devCmd)
	dev.SetInteractiveCommand(devCmd, dev.DevInteractive)
}

func devRun(cmd *cobra.Command, args []string) error {
	output, _ := cmd.Flags().GetString("output")
	return dev.RunGo(args, output)
}

func devRunAdd(cmd *cobra.Command, args []string) error {
	msl, _ := cmd.Flags().GetUint8("max-sublevels")
	return dev.RunAddFolder(args[0], msl)
}

func devRunRemove(cmd *cobra.Command, args []string) error {
	return dev.RunRemoveFolder(args[0])
}

func devRunSync(cmd *cobra.Command, args []string) error {
	return dev.RunSync()
}

func devRunChild(cmd *cobra.Command, args []string) error {
	enable, _ := cmd.Flags().GetBool("enable")
	disable, _ := cmd.Flags().GetBool("disable")
	if !(enable || disable) {
		return errors.New("you must inform --enable or --disable")
	}

	folder, _ := os.Getwd()

	if len(args) > 0 {
		folder = args[0]
	}
	return dev.RunIgnoreChildren(folder, disable)

}
