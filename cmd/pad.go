package cmd

import (
	"fmt"
	"time"

	"github.com/guionardo/gs-dev/app/pad"
	"github.com/spf13/cobra"
)

// goCmd represents the go command
var padCmd = &cobra.Command{
	Use:   "pad",
	Short: "Share text by URL",
	Long:  `Share your text/data/code using a external pad tool`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("pad called")
	},
}

func init() {
	setupCmd := &cobra.Command{
		Use:  "setup",
		RunE: pad.RunSetupPad,
	}
	setupCmd.Flags().String("backend", "https://gs-bucket.fly.dev", "Bucket storage URL (https)")
	setupCmd.Flags().String("api-key", "", "API Key to publish pads")

	sendFileCmd := &cobra.Command{
		Use:   "file",
		Short: "File name to send",
		RunE:  pad.RunSendFilePad,
		Args:  cobra.ExactArgs(1),
	}

	sendTextCmd := &cobra.Command{
		Use:   "text",
		Short: "Text to send (use quotes)",
		Long:  "Args: <file name> \"content\"",
		RunE:  pad.RunSendTextPad,
		Args:  cobra.ExactArgs(2),
	}

	for _, subCmd := range []*cobra.Command{sendFileCmd, sendTextCmd} {
		subCmd.Flags().String("slug", "", "Easy to remember name for URL")
		subCmd.Flags().Duration("ttl", time.Duration(0), "Time to live")
		subCmd.Flags().Bool("delete-after-read", false, "Delete pad after first read/download")
	}
	padCmd.AddCommand(
		setupCmd,
		sendFileCmd,
		sendTextCmd,
	)

	rootCmd.AddCommand(padCmd)

}
