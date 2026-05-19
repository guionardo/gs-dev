package padcommand

import (
	"fmt"
	"time"

	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/dialog"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"
	"github.com/spf13/cobra"
)

type PadCliCommand struct {
	commands.CommonCommand

	configuration *config.ConfigRoot
}

const (
	cliName        = "pad"
	cliDescription = "Pad is a form of transmit short information."

	getCommand     = "get"
	getDescription = "get a pad"

	getLastCommand     = "last"
	getLastDescription = "get last posted pads"
)

func (p *PadCliCommand) Init() {
	p.InitCommandVariables(cliName, cliDescription, p, p.setup).WithInitAlias().WithOutput()
}

func (p *PadCliCommand) setup(configuration *config.ConfigRoot) (err error) {
	p.configuration = configuration
	return nil
}

func (p *PadCliCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&PadCliStruct{}, cliName, cliDescription, "", false)
}

func (p *PadCliCommand) GetTUICommand() func() error {
	return func() error {
		service := padservice.NewPadClientService(p.configuration)

		command, err := dialog.Choose(
			"dev",
			dialog.ToAnyArray([]string{
				fmt.Sprintf("%s:%s", getCommand, getDescription),
				fmt.Sprintf("%s:%s", getLastCommand, getLastDescription),
			})...)
		if err != nil {
			return err
		}

		switch command {
		case getCommand:
			postId, err := dialog.Input("Post ID")
			if err != nil {
				return err
			}

			content, err := service.Get(postId)
			if err == nil {
				fmt.Print(string(content))
			}

			return err
		case getLastCommand:
			options := make([]string, 0)
			for _, post := range service.GetLastPostIDs() {
				options = append(options, fmt.Sprintf("%s:%s", post.ID, post.CreatedAt.Format(time.DateTime)))
			}

			if len(options) == 0 {
				fmt.Print("There no last posts registered")
				return nil
			}

			postId, err := dialog.Choose("last pads", dialog.ToAnyArray(options)...)
			if err != nil {
				return err
			}

			content, err := service.Get(postId)
			if err == nil {
				fmt.Print(string(content))
			}

			return err
		}

		// fmt.Println(description)

		return nil
	}
}
