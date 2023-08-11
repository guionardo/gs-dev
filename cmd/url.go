package cmd

import (
	"os"

	"github.com/guionardo/gs-dev/app/url"
	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "Open the git repository on browser",
	RunE:  urlRun,
}

func init() {
	urlCmd.Flags().BoolP("just-show", "j", false, "Doesn't open, just print the URL")
	rootCmd.AddCommand(urlCmd)
}

func urlRun(cmd *cobra.Command, args []string) (err error) {
	var folder string
	if len(args) == 0 {
		folder, err = os.Getwd()
	} else {
		folder = args[0]
	}
	if err != nil {
		return
	}
	if justShow, err := cmd.Flags().GetBool("just-show"); err != nil {
		return err
	} else {
		return url.RunUrl(folder, justShow)
	}
}
