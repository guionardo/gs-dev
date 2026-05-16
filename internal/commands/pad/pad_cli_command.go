package padcommand

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/guionardo/gs-dev/internal/cli"
	"github.com/guionardo/gs-dev/internal/commands"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/configurations"
	errs "github.com/guionardo/gs-dev/internal/errors"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"
	string_tools "github.com/guionardo/gs-dev/internal/tools/strings_tools"
	"github.com/spf13/cobra"
)

type PadCliCommand struct {
	commands.CommonCommand

	service *padservice.PadClientService
}

const (
	cliName        = "pad"
	cliDescription = "Pad is a form of transmit short information."
)

func (p *PadCliCommand) Init() {
	p.InitCommandVariables(cliName, cliDescription, p, p.setup).WithInitAlias().WithOutput()
}

func (p *PadCliCommand) setup(configuration *config.ConfigRoot) (err error) {
	p.service = padservice.NewPadClientService(configuration)

	return nil
}

func (p *PadCliCommand) BuildCommand() *cobra.Command {
	return cli.GenerateCobraCommand(&PadCliStruct{}, cliName, cliDescription, "", false)
}

func (p *PadCliCommand) _BuildCommand() *cobra.Command {
	clientConfig := p.service.GetConfig()
	cmd := &cobra.Command{
		Use:   cliName,
		Short: cliDescription,
	}
	cmd.PersistentFlags().String("host", clientConfig.BackendURL, "Pad server backend URL")
	cmd.PersistentFlags().String("api-key", clientConfig.APIKey, "Pad server API key")
	cmd.AddCommand(p.getCmdGet())
	cmd.AddCommand(p.getCmdPost())
	cmd.AddCommand(p.getCmdDel())

	return cmd
}

func (p *PadCliCommand) GetTUICommand() func() error {
	return func() error {
		return nil //TODO: Implement
	}
}

func (p *PadCliCommand) getCmdGet() *cobra.Command {
	cmdGet := &cobra.Command{
		Use:   "get",
		Short: "Get a pad by ID",
		RunE:  p.runGetPad,
		Args:  cobra.ExactArgs(1),
	}

	return cmdGet
}

func (p *PadCliCommand) getConfigFunc(cmd *cobra.Command) {
	var (
		config     *configurations.PadClientConfig
		backendURL string
		apiKey     string
	)

	config = p.service.GetConfig()

	if cmd.Flags().Changed("host") {
		backendURL, _ = cmd.Flags().GetString("host")
	}

	if cmd.Flags().Changed("api-key") {
		apiKey, _ = cmd.Flags().GetString("api-key")
	}

	if len(backendURL) > 0 || len(apiKey) > 0 {
		newConfig := *config
		newConfig.BackendURL = backendURL
		newConfig.APIKey = apiKey
		newConfig.Enabled = true
		p.service.SetCustomClientConfig(&newConfig)
	} else {
		p.service.SetCustomClientConfig(nil)
	}
}

func (p *PadCliCommand) getCmdPost() *cobra.Command {
	cmdPost := &cobra.Command{
		Use:   "post",
		Short: "Create a new pad",
		RunE:  p.runPostPad,
	}
	cmdPost.Flags().DurationP("ttl", "t", 0, "Time to live for the pad (e.g., 10m, 1h). 0 means no expiration.")
	cmdPost.Flags().StringSlice("header", []string{}, "Additional headers to include in the request (format: Key=Value)")
	cmdPost.Flags().StringP("input", "i", "", "Input content for the pad. If not provided, content will be read from stdin.")

	return cmdPost
}

func (p *PadCliCommand) getCmdDel() *cobra.Command {
	cmdDel := &cobra.Command{
		Use:   "del",
		Short: "Delete a pad by ID",
		RunE:  p.runDelPad,
		Args:  cobra.ExactArgs(1),
	}

	return cmdDel
}

func (p *PadCliCommand) runGetPad(cmd *cobra.Command, args []string) error {
	postID := args[0]

	p.getConfigFunc(cmd)

	content, err := p.service.Get(postID)
	if err != nil {
		return err
	}

	_, err = cmd.OutOrStdout().Write(content)

	return err
}

func (p *PadCliCommand) runPostPad(cmd *cobra.Command, args []string) error {
	var (
		ttl       time.Duration
		content   []byte
		headerMap map[string]string
		err       error
	)
	// Get TLL from flags
	if cmd.Flags().Changed("ttl") {
		ttl, err = cmd.Flags().GetDuration("ttl")
		if err != nil {
			return err
		}
	}
	// Get headers from flags
	if cmd.Flags().Changed("header") {
		headers, err := cmd.Flags().GetStringSlice("header")
		if err != nil {
			return err
		}

		headerMap = make(map[string]string)

		for _, header := range headers {
			var key, value string
			string_tools.SplitString(header, "=", &key, &value)

			if key == "" || value == "" {
				return fmt.Errorf("invalid header format: %s", header)
			}

			headerMap[key] = value
		}
	}

	// Get input content

	input, err := cmd.Flags().GetString("input")

	var inputReader io.Reader
	if input == "" {
		inputReader = cmd.InOrStdin()
	} else {
		file, fileErr := os.Open(input)

		err = fileErr
		if fileErr != nil {
			return errs.NewError(fileErr, "error opening input file", false)
		}

		defer file.Close()

		inputReader = file
	}

	content, err = io.ReadAll(inputReader)
	if err != nil {
		return errs.NewError(err, "error reading input content: %s", false, err.Error())
	}

	postID, err := p.service.Post(content, ttl, headerMap)
	if err != nil {
		return errs.NewError(err, "error posting pad content", false)
	}

	cmd.Printf("Pad uploaded. ID = %s", postID)

	return nil
}

func (p *PadCliCommand) runDelPad(cmd *cobra.Command, args []string) error {
	postID := args[0]

	p.getConfigFunc(cmd)

	err := p.service.Delete(postID)
	if err != nil {
		return err
	}

	cmd.Printf("Pad with ID %s deleted successfully.", postID)

	return nil
}
