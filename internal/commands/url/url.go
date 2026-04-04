package url

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	urlservice "github.com/guionardo/gs-dev/internal/services/url"
	"github.com/spf13/cobra"
)

type UrlCommand struct {
	commands.CommonCommand

	service *urlservice.URLService
}

const (
	name        = "url"
	description = "Open a URL"
)

func (u *UrlCommand) Init() {
	u.InitCommandVariables(name, description, false, "", false)
}

func (u *UrlCommand) Setup(configuration *config.ConfigFile) error {
	u.service = urlservice.NewURLService()
	return nil
}

func (u *UrlCommand) GetCobraCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			repositoryRoot, err := cmd.Flags().GetString("directory")
			if err != nil {
				return err
			}

			if justShow, err := cmd.Flags().GetBool("just-show"); err == nil && justShow {
				url, err := u.service.GetRemoteHttpURL(repositoryRoot)
				if err != nil {
					return err
				}

				cmd.Print(url)

				return nil
			}

			return u.service.OpenRemoteURL(repositoryRoot)
		},
	}
	cmd.Flags().StringP("directory", "d", ".", "Directory path")
	cmd.Flags().BoolP("just-show", "j", false, "Just show and doesn´t open in browser")

	return cmd
}

func (u *UrlCommand) GetTUICommand() func() error {
	return func() error {
		url, err := u.service.GetRemoteHttpURL(".")
		if err != nil {
			return err
		}

		if dialog.Confirm(fmt.Sprintf("Open URL [%s] in browser for %s?", url, u.GetRepositoryName()), true) {
			return u.service.OpenRemoteURL(".")
		}

		return nil
	}
}
