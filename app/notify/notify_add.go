package notify

import (
	"fmt"
	"strings"

	"github.com/guionardo/gs-dev/configs"
	"github.com/spf13/cobra"
)

func AddNotify(cmd *cobra.Command, args []string) error {
	message := strings.Join(args, " ")

	if len(message) == 0 {
		return fmt.Errorf("empty message")
	}
	messageType, err := cmd.Flags().GetString("type")
	if err != nil {
		return err
	}
	var mt int
	switch strings.ToLower(messageType) {
	case "warning":
		mt = MESSAGE_WARNING
	case "error":
		mt = MESSAGE_ERROR
	default:
		mt = MESSAGE_INFO
	}
	cfg, err := configs.GetConfiguration(cmd)
	if err != nil {
		return err
	}

	repository := CreateNotifyRepository(cfg.DataFolder)
	_, err = repository.AddNotification(message, mt)
	return err
}
