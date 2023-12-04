package config

type Reader interface {
	Read() (Config, error)
}
