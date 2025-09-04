package config

import "github.com/patrickhuber/go-config"

type DefaultProvider interface {
	GetDefault(ctx *config.GetContext) (Config, error)
}
