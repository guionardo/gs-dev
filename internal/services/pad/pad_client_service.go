package padservice

import (
	"time"

	"github.com/guionardo/go/flow"
	"github.com/guionardo/gs-dev/internal/compression"
	"github.com/guionardo/gs-dev/internal/config"
	pad_client "github.com/guionardo/gs-dev/internal/pad/client"
)

type PadClientService struct {
	configRoot   *config.ConfigRoot
	config       *PadClientConfig
	customConfig *PadClientConfig
	client       *pad_client.PadClient
}

const (
	ContentCompressorHeader = "X-Content-Compressor"
)

func NewPadClientService(configuration *config.ConfigRoot) *PadClientService {
	clientConfig, err := config.GetValue[PadClientConfig](configuration)
	if err != nil {
		clientConfig = PadClientConfig{}
	}

	clientConfig.Defaults()
	service := &PadClientService{configRoot: configuration, config: &clientConfig}

	if clientConfig.BackendURL != "" {
		service.client, _ = pad_client.NewPadClient(clientConfig.BackendURL, clientConfig.APIKey)
	}

	return service
}

func (p *PadClientService) Post(content []byte, ttl time.Duration, headers map[string]string) (postID string, err error) {
	// Try to compress the content before sending it to the server
	compressed, compressor := compression.Compress(content)

	if headers == nil {
		headers = make(map[string]string)
	}

	headers[ContentCompressorHeader] = compressor
	validUntil := flow.If(ttl > 0, time.Now().Add(ttl), time.Time{}).Round(time.Second)

	postID, err = p.client.Post(compressed, ttl, headers)
	if err == nil {
		p.config.LastPosts.Add(postID, validUntil)
		config.SetValue(p.configRoot, p.config)
		_ = p.configRoot.Save()
	}

	return postID, err // TODO: Wrap error with more context
}

func (p *PadClientService) Get(postID string) (content []byte, err error) {
	content, metadata, err := p.client.Get(postID)

	if err == nil && metadata != nil {
		content, err = compression.Decompress(content, metadata[ContentCompressorHeader])
	}

	if err == nil {
		p.config.LastPosts.Add(postID, time.Now().AddDate(0, 0, 7))
	} else {
		p.config.LastPosts.Remove(postID)
	}

	config.SetValue(p.configRoot, p.config)
	_ = p.configRoot.Save()

	return content, err // TODO: Wrap error with more context
}

func (p *PadClientService) Delete(postID string) (err error) {
	err = p.client.Delete(postID)

	return err
}

func (p *PadClientService) GetConfig() *PadClientConfig {
	if p.customConfig != nil {
		return p.customConfig
	}

	return p.config
}

func (p *PadClientService) SaveConfig(cfg *PadClientConfig) error {
	config.SetValue(p.configRoot, cfg)

	if err := cfg.Validate(); err != nil {
		return err
	}

	return p.configRoot.Save()
}

func (p *PadClientService) SetCustomClientConfig(clientConfig *PadClientConfig) {
	p.customConfig = clientConfig
}

func (p *PadClientService) GetLastPostIDs() []Post {
	return p.config.LastPosts.Posts()
}
