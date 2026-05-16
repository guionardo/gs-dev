package configurations

type StorageConfig struct {
	Enabled bool              `yaml:"enabled"`
	Type    string            `yaml:"type"`
	Options map[string]string `yaml:"options"`
}

func (sc *StorageConfig) Defaults() {
	if sc.Type == "" {
		sc.Type = "fs" // file system
	}

	if sc.Options == nil {
		sc.Options = make(map[string]string, 0)
	}
}

func (sc StorageConfig) Key() string {
	return "store"
}

func (sc StorageConfig) Validate() error {
	return nil // TODO: Implementar validação
}
