package storage

import (
	"context"
	"log/slog"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/interfaces"
)

// NewStorage creates a Storage from the configuration. Currently only
// FileSystemStorage is supported; other implementations will be added in
// the future.
func NewStorage(configuration *config.ConfigRoot, ctx context.Context) (interfaces.Storage, error) {
	// TODO: Implement another types of storage
	return NewFileSystemStorage(configuration, ctx, slog.Default())
}
