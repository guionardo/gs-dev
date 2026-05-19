package url

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/cli"
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
	u.InitCommandVariables(name, description, u, u.setup)
}

func (u *UrlCommand) setup(configuration *config.ConfigRoot) error {
	u.service = urlservice.NewURLService()
	return nil
}

func (u *UrlCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&UrlStruct{}, name, description, "", false)
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
