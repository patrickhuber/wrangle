package config

type Writer interface {
	Write(c *Config) error
}
