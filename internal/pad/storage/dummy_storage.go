package storage

import (
	"time"

	"github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/interfaces"
)

// DummyStorage implements Storage by returning disabledStorageError for every
// operation. It is used when storage is disabled in configuration.
type DummyStorage struct {
}

var (
	_                  interfaces.PadStorage = &DummyStorage{}
	errDisabledStorage                       = errors.NewError(nil, "storage is disabled", false)
)

func (ds *DummyStorage) Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error) {
	return "", errDisabledStorage
}
func (ds *DummyStorage) Get(postID string) (content []byte, headers map[string]string, err error) {
	return nil, nil, errDisabledStorage
}
func (ds *DummyStorage) Delete(postID string) error {
	return errDisabledStorage
}
