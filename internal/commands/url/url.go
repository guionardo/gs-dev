package url

import (
	"fmt"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/spf13/cobra"
)

type UrlCommand struct {
}

const (
	name        = "url"
	description = "Open a URL"
)

func (u *UrlCommand) Setup(configuration *config.ConfigFile) error {
	return nil
}

func (u *UrlCommand) GetName() string {
	return name
}

func (u *UrlCommand) GetDescription() string {
	return description
}

func (u *UrlCommand) GetCobraCommand() *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: description,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(description)
			return nil
		},
	}
}

func (u *UrlCommand) GetTUICommand() func() error {
	return func() error {
		fmt.Println(description)
		return nil
	}
}
