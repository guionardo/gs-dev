package padservice

import (
	"time"

	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/errors"
	"github.com/guionardo/gs-dev/internal/interfaces"
)

type PadServerService struct {
	config       *PadServerConfig
	customConfig *PadServerConfig
	storage      interfaces.PadStorage
}

func NewPadServerService(configuration *config.ConfigRoot, storage interfaces.PadStorage) (*PadServerService, error) {
	configServer, _ := config.GetValue[PadServerConfig](configuration)
	if err := configServer.Validate(); err != nil {
		return nil, err
	}

	if storage == nil {
		return nil, errors.NewError(nil, "storage can not be nil", false)
	}

	return &PadServerService{
		config:  &configServer,
		storage: storage,
	}, nil
}

func (p *PadServerService) Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error) {
	return p.storage.Post(content, ttl, headers)
}

func (p *PadServerService) Get(postID string) (content []byte, headers map[string]string, err error) {
	return p.storage.Get(postID)
}

func (p *PadServerService) Delete(postID string) error {
	return p.storage.Delete(postID)
}

func (p *PadServerService) IsAPIKeyValid(apiKey string) bool {
	return p.isAuthorized(apiKey)
}

func (p *PadServerService) GetDefaultTTL() time.Duration {
	return min(p.config.StorageConfig.DefaultTTL, time.Duration(0))
}

func (p *PadServerService) isAuthorized(apiKey string) bool {
	return len(p.config.APIKey) == 0 || p.config.APIKey == apiKey
}

func (p *PadServerService) GetConfig() *PadServerConfig {
	if p.customConfig != nil {
		return p.customConfig
	}

	return p.config
}
func (p *PadServerService) SetCustomClientConfig(serverConfig *PadServerConfig) {
	p.customConfig = serverConfig
}
