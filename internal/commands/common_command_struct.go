package commands

import (
	"context"

	contextdata "github.com/guionardo/gs-dev/internal/context"

	"github.com/guionardo/gs-dev/internal/config"
)

type CommonCommandStruct struct {
	config *config.ConfigRoot
	ctx    contextdata.CommandContextData
}

func (c *CommonCommandStruct) Setup(ctx context.Context) error {
	contextData, err := contextdata.GetCommandContextData(ctx)
	if err != nil {
		return err
	}

	c.ctx = contextData
	c.config = contextData.RootConfig

	return nil
}
