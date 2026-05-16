package padcommand

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/guionardo/gs-dev/internal/cli"
	ctxdata "github.com/guionardo/gs-dev/internal/context"
	"github.com/guionardo/gs-dev/internal/errors"
	padservice "github.com/guionardo/gs-dev/internal/services/pad"
)

type (
	PadCliStruct struct {
		basePadCliStruct

		Host   string `flag:"host" description:"Pad server backend URL" persistent:"true"`
		ApiKey string `flag:"api-key" description:"Pad server API key" persistent:"true"`

		Get     PadCliGetStruct   `subcommand:"get"`
		Post    PadCliPostStruct  `subcommand:"post"`
		Del     PadCliDelStruct   `subcommand:"del"`
		DoSetup PadCliSetupStruct `subcommand:"setup"`
	}

	PadCliGetStruct struct {
		basePadCliStruct

		Output string `flag:"output,o" description:"Output for a file or stdout (default)" default:"stdout"`

		Args []string `args:"1,Post ID Example ABCD-1234"`
	}
	PadCliPostStruct struct {
		basePadCliStruct

		Filename string        `flag:"filename,f" description:"existing file to create a pad."`
		StdIn    bool          `flag:"stdin" description:"Read data from std in (pipe)`
		TTL      time.Duration `flag:"ttl" description:"time to live of the pad"`
		Headers  []string      `flag:"header" description:"header=value"`
		// Body     []byte        `stdin:"true"`
	}
	PadCliDelStruct struct {
		basePadCliStruct

		Args []string `args:"1,Post ID Example ABCD-1234"`
	}

	PadCliSetupStruct struct {
		basePadCliStruct

		Show       bool   `flag:"show" default:"false"`
		Enabled    bool   `flag:"enabled" default:"true"`
		BackendURL string `flag:"backend_url" default:"http://localhost:8080"`
		APIKey     string `flag:"api_key" default:""`
	}

	basePadCliStruct struct {
		service *padservice.PadClientService
	}
)

func (i *PadCliStruct) Run(ctx context.Context, output io.Writer) error {
	return errors.NewError(nil, "a subcommand is required", false)
}
func (i *basePadCliStruct) Setup(ctx context.Context) error {
	cd, err := ctxdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	i.service = padservice.NewPadClientService(cd.RootConfig)

	return nil
}

func (cg *PadCliGetStruct) Run(ctx context.Context, output io.Writer) error {
	postId := cg.Args[0]

	content, err := cg.service.Get(postId)
	if err != nil {
		return err
	}

	if cg.Output == "stdout" {
		_, err = output.Write(content)
		return err
	}

	file, err := os.Create(cg.Output)
	if err == nil {
		defer file.Close()

		_, err = file.Write(content)
	}

	return err
}

func (cp *PadCliPostStruct) Run(ctx context.Context, output io.Writer) (err error) {
	var (
		content []byte
		headers = map[string]string{}
	)

	if cp.StdIn {
		content, err = cli.ReadFromStdIn()
	} else if cp.Filename != "" {
		file, err := os.Open(cp.Filename)
		if err != nil {
			return errors.NewError(err, "failed to open file: %s", false, cp.Filename)
		}
		defer file.Close()

		content, err = io.ReadAll(file)
		headers["filename"] = path.Base(cp.Filename)
	} else {
		err = errors.NewError(nil, "required source of data for pad", false)
	}

	if err != nil {
		return err
	}

	if len(content) == 0 {
		return errors.NewError(nil, "cannot post empty pad", false)
	}

	mimeType := http.DetectContentType(content)
	headers["mime-type"] = mimeType

	for _, header := range cp.Headers {
		w := strings.SplitN(header, "=", 2)
		if len(w) == 2 {
			headers[w[0]] = w[1]
		}
	}

	postID, err := cp.service.Post(content, cp.TTL, headers)
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "Created new pad: %s\n\n", postID)

	for h, v := range headers {
		fmt.Fprintf(output, "%s=%v\n", h, v)
	}

	return nil
}

func (cp *PadCliDelStruct) Run(ctx context.Context, output io.Writer) error {
	postID := cp.Args[0]

	err := cp.service.Delete(postID)
	if err != nil {
		return err
	}

	fmt.Fprintf(output, "Pad deleted: %s\n", postID)

	return nil
}

func (cp *PadCliSetupStruct) Run(ctx context.Context, output io.Writer) error {
	cfg := cp.service.GetConfig()

	if cp.Show {
		fmt.Fprintf(output, `PAD CLI SETUP
		
Enabled = %v
Backend URL = %s
API Key = %s`, cfg.Enabled, cfg.BackendURL, cfg.APIKey)

		return nil
	}

	cfg.APIKey = cp.APIKey
	cfg.BackendURL = cp.BackendURL
	cfg.Enabled = cp.Enabled

	return cp.service.SaveConfig(cfg)
}
