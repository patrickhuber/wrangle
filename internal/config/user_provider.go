package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/global"
)

type userProvider struct {
	fs fs.FS
}

func NewUserProvider(fs fs.FS) config.Provider {
	return &userProvider{
		fs: fs,
	}
}

func (p *userProvider) Get(ctx *config.GetContext) (any, error) {
	// use a dataptr to get the value of the env config
	userConfigPath, err := dataptr.GetAs[string]("/spec/env/"+global.EnvUserConfig, ctx.MergedConfiguration)
	if err != nil {
		return nil, err
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
