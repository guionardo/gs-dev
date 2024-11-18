package config

type Configuration interface {
	Load() error
	Save() error
}
