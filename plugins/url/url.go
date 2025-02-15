package url

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/guionardo/gs-dev/internal/git"
	"github.com/guionardo/gs-dev/pkg/plugins"
	"github.com/guionardo/gs-dev/plugins/commons"
	"github.com/spf13/cobra"
)

type URL struct {
	commons.BasePlugin
}

const urlName = "url"

var (
	directory string
	justShow  *bool
)

func NewURL() plugins.CliPlugin {
	return &URL{
		BasePlugin: commons.BasePlugin{
			PluginName:    urlName,
			CanBeDisabled: true,
			Enabled:       true,
		},
	}
}

func (u *URL) RunURL(command *cobra.Command, args []string) (err error) {
	slog.Debug("Getting URL from git project", slog.String("directory", directory), slog.Bool("justShow", *justShow))
	var url string
	if url, err = git.GetRemoteHttpURL(directory); err == nil {
		fmt.Print(url)
		if !*justShow {
			err = openInBrowser(url)
		}
	}
	return

}
func (u *URL) GetCobraCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   urlName,
		Short: "Get URL for repository",
		RunE:  u.RunURL,
	}
	wd, err := os.Getwd()
	if err != nil {
		wd = "."
	}
	if awd, err := filepath.Abs(wd); err == nil {
		wd = awd
	}
	cmd.Flags().StringVarP(&directory, "directory", "d", wd, "Directory path")
	// cmd.Flags().Lookup("directory").NoOptDefVal = wd

	justShow = cmd.Flags().BoolP("just-show", "j", false, "Just show and doesn´t open in browser")

	return cmd
}

func (u *URL) Setup(manager plugins.PluginsManager, configurationFolder string) error {
	return nil
}

func checkReachableUrl(url string) error {
	if resp, err := http.Head(url); err != nil {
		return err
	} else if resp.StatusCode >= 100 {
		return nil
	} else {
		return fmt.Errorf("got %s status from %s", resp.Status, url)
	}
}

func openInBrowser(url string) (err error) {
	if err = checkReachableUrl(url); err != nil {
		return
	}
	command, args := urlCommand(url)
	err = exec.Command(command, args...).Start()

	return
}
