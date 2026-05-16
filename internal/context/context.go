package context

import (
	"context"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/errors"
)

type (
	// CommandContextData holds the data that will be passed to the command execution context
	CommandContextData struct {
		RootConfig *config.ConfigRoot
	}

	contextKey string
)

const CommandContextKey contextKey = "command"

// GetCommandContext creates a new context with the provided CommandContextData
func GetCommandContext(originalContext context.Context, ctxData CommandContextData) context.Context {
	return context.WithValue(originalContext, CommandContextKey, ctxData)
}

// GetCommandContextData retrieves the CommandContextData from the context
func GetCommandContextData(ctx context.Context) (CommandContextData, error) {
	data, ok := ctx.Value(CommandContextKey).(CommandContextData)
	if !ok {
		return data, errors.ErrCommandContextDataNotFound
	}

	return data, nil
}
