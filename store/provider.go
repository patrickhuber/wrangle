package store

import "github.com/patrickhuber/wrangle/config"

// Provider provides a store given the config source
type Provider interface {
	GetName() string
	Create(source *config.ConfigSource) (Store, error)
}
