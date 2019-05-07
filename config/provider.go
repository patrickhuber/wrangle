package config

// Provider defines methods for getting and setting config
type Provider interface {
	Get() (*Config, error)
	Set(config *Config) error
}
