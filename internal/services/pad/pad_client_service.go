package padservice

import (
	"time"

	"github.com/guionardo/gs-dev/internal/compression"
	"github.com/guionardo/gs-dev/internal/config"
	"github.com/guionardo/gs-dev/internal/configurations"
	pad_client "github.com/guionardo/gs-dev/internal/pad/client"
)

type PadClientService struct {
	configRoot   *config.ConfigRoot
	config       *configurations.PadClientConfig
	customConfig *configurations.PadClientConfig
	client       *pad_client.PadClient
}

const (
	ContentCompressorHeader = "X-Content-Compressor"
)

func NewPadClientService(configuration *config.ConfigRoot) *PadClientService {
	clientConfig, err := config.GetValue[configurations.PadClientConfig](configuration)
	if err != nil {
		clientConfig = configurations.PadClientConfig{}
	}

	clientConfig.Defaults()
	service := &PadClientService{configRoot: configuration, config: &clientConfig}

	var client *pad_client.PadClient
	if clientConfig.BackendURL != "" {
		service.client, _ = pad_client.NewPadClient(clientConfig.BackendURL, clientConfig.APIKey)
	}

	return &PadClientService{
		configRoot: configuration,
		config:     &clientConfig,
		client:     client,
	}
}

func (p *PadClientService) Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error) {
	// Try to compress the content before sending it to the server
	compressed, compressor := compression.Compress(content)

	if headers == nil {
		headers = make(map[string]string)
	}

	headers[ContentCompressorHeader] = compressor
	postID, err = p.client.Post(compressed, ttl, headers)

	return postID, err // TODO: Wrap error with more context
}

func (p *PadClientService) Get(postID string) (content []byte, err error) {
	content, metadata, err := p.client.Get(postID)

	if err == nil && metadata != nil {
		content, err = compression.Decompress(content, metadata[ContentCompressorHeader])
	}

	return content, err // TODO: Wrap error with more context
}

func (p *PadClientService) Delete(postID string) (err error) {
	err = p.client.Delete(postID)

	return err
}

func (p *PadClientService) GetConfig() *configurations.PadClientConfig {
	if p.customConfig != nil {
		return p.customConfig
	}

	return p.config
}

func (p *PadClientService) SaveConfig(cfg *configurations.PadClientConfig) error {
	config.SetValue(p.configRoot, cfg)

	if err := cfg.Validate(); err != nil {
		return err
	}

	return p.configRoot.Save()
}

func (p *PadClientService) SetCustomClientConfig(clientConfig *configurations.PadClientConfig) {
	p.customConfig = clientConfig
}
