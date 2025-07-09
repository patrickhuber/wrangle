package config

import (
	"fmt"

	iofs "io/fs"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/global"
)

type userProvider struct {
	fs               fs.FS
	errorIfNotExists bool
}

func NewUserProvider(fs fs.FS, errorIfNotExists bool) config.Provider {
	return &userProvider{
		fs:               fs,
		errorIfNotExists: errorIfNotExists,
	}
}

func (p *userProvider) Get(ctx *config.GetContext) (any, error) {
	// use a dataptr to get the value of the env config
	userConfigPath, err := dataptr.GetAs[string]("/spec/env/"+global.EnvUserConfig, ctx.MergedConfiguration)
	if err != nil {
		return nil, err
	}

	// if the user config doesn't exist and the errorIfNotExists flag is set, return an error
	if p.errorIfNotExists {
		exists, err := p.fs.Exists(userConfigPath)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w user config file %s does not exist", iofs.ErrNotExist, userConfigPath)
		}
	}

	// does the file exist, if not, create it
	cfg, err := ReadOrCreateFile(p.fs, userConfigPath, p.getDefault)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (p *userProvider) getDefault() Config {
	return Config{}
}
