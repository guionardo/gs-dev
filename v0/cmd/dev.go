package cmd

import (
	"errors"
	"os"
	"path"
	"time"

	"github.com/guionardo/gs-dev/app"
	"github.com/guionardo/gs-dev/app/dev"
	"github.com/spf13/cobra"
)

const (
	maxSubLevels    = "max-sub-levels"
	Enable          = "enable"
	Disable         = "disable"
	MaxSyncInterval = "max-sync-interval"
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
	addCommand.Flags().Uint8P(maxSubLevels, "m", 3, "Maximum sub level for folders")

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
	syncCommand.Flags().Int16P(MaxSyncInterval, "m", 60, "Max synchronization interval in minutes")

	ignoreChildCommand := &cobra.Command{
		Use:     "child",
		Short:   "Enable/disable children",
		Long:    "use arguments 'enable' or 'disable'",
		Args:    cobra.ExactArgs(1),
		RunE:    devRunChild,
		Example: "dev child --enable /home/some/folder",
	}
	ignoreChildCommand.Flags().BoolP(Enable, "e", false, "Enable children")
	ignoreChildCommand.Flags().BoolP(Disable, "d", false, "Disable children")

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
	msl, _ := cmd.Flags().GetUint8(maxSubLevels)
	return dev.RunAddFolder(args[0], msl)
}

func devRunRemove(cmd *cobra.Command, args []string) error {
	return dev.RunRemoveFolder(args[0])
}

func devRunSync(cmd *cobra.Command, args []string) error {
	maxSyncInterval, _ := cmd.Flags().GetInt16(MaxSyncInterval)
	return dev.RunSync(time.Minute * time.Duration(maxSyncInterval))
}

func devRunChild(cmd *cobra.Command, args []string) error {
	enable, _ := cmd.Flags().GetBool(Enable)
	disable, _ := cmd.Flags().GetBool(Disable)
	if !(enable || disable) {
		return errors.New("you must inform --enable or --disable")
	}

	folder, _ := os.Getwd()

	if len(args) > 0 {
		folder = args[0]
	}
	return dev.RunIgnoreChildren(folder, disable)

}
