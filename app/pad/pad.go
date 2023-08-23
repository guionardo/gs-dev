package pad

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	bucket "github.com/guionardo/gs-bucket/client"
	"github.com/guionardo/gs-dev/config"
	"github.com/spf13/cobra"
)

func getPadClient(cfg *config.PadConfig) (c *bucket.GSBucketClient, err error) {
	if cfg == nil || len(cfg.BackendUrl) == 0 {
		err = errors.New("Missing backend URL configuration")
		return
	}
	c = bucket.CreateGSBucketClient(cfg.BackendUrl, cfg.ApiKey)
	return
}

func RunSetupPad(cmd *cobra.Command, args []string) (err error) {
	backend, _ := cmd.Flags().GetString("backend")
	apiKey, _ := cmd.Flags().GetString("api-key")
	config := getConfig()
	changed := false
	if len(backend) > 0 {
		if _, err = url.ParseRequestURI(backend); err != nil {
			return err
		}
		config.BackendUrl = backend
		changed = true
	}
	if len(apiKey) > 0 {
		config.ApiKey = apiKey
		changed = true
	}
	if changed {
		err = config.Save()
	}

	return
}

func RunSendFilePad(cmd *cobra.Command, args []string) error {
	fileName := args[0]
	if content, err := os.ReadFile(fileName); err != nil {
		return err
	} else {
		return runSendPad(cmd, fileName, content)
	}
}

func RunSendTextPad(cmd *cobra.Command, args []string) error {
	fileName := args[0]
	content := []byte(args[1])
	return runSendPad(cmd, fileName, content)
}

func runSendPad(cmd *cobra.Command, fileName string, content []byte) error {
	config := getConfig()
	client, err := getPadClient(config)
	if err != nil {
		return err
	}

	ttl, _ := cmd.Flags().GetDuration("ttl")
	deleteAfterRead, _ := cmd.Flags().GetBool("delete-after-read")
	slug, _ := cmd.Flags().GetString("slug")

	pad, err := client.CreatePad(fileName, ttl, deleteAfterRead, content, slug)
	if err != nil {
		return err
	}
	fmt.Printf("Pad created: %s/pads/%s", config.BackendUrl, pad.Slug)
	return nil
}
