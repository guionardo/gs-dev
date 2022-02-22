package configs

type RootConfig struct {
	DataFolder string    `yaml:"data_folder"`
	DevConfig  DevConfig `yaml:"dev_config"`
}
