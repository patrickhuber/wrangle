package store

import "github.com/patrickhuber/wrangle/config"

// Provider provides a store given the config source
type Provider interface {
	Name() string
	Create(source *config.Store) (Store, error)
}
