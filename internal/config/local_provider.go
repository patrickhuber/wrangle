package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
)

type localProvider struct {
	filesystem       fs.FS
	path             filepath.Provider
	errorIfNotExists bool
}

func NewLocalProvider(fs fs.FS, path filepath.Provider, errorIfNotExists bool) config.Provider {
	// use the config.NewGlobUp function to load the local configuration

	return &localProvider{
		filesystem:       fs,
		path:             path,
		errorIfNotExists: errorIfNotExists,
	}
}

func (p *localProvider) Get(ctx *config.GetContext) (any, error) {

	return nil, nil
}

func (p *localProvider) getDefault() Config {
	return Config{}
}
