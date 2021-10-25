package config

// Reader defines methods for reading config
type Reader interface {
	Get() (*Config, error)
}

type Writer interface {
	Set(config *Config) error
}

// Provider defines methods for getting and setting config
type Provider interface {
	Reader
	Writer
}
