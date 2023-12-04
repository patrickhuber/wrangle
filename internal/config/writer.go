package config

type Writer interface {
	Write(Config) error
}
