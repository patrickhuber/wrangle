package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/global"
)

func NewSystemProvider(fs fs.FS) config.Provider {
	return &systemProvider{
		fs: fs,
	}
}

type systemProvider struct {
	fs fs.FS
}

// Get implements config.Provider.
func (p *systemProvider) Get(ctx *config.GetContext) (any, error) {
	// use a dataptr to get the value of the env config
	systemConfigPath, err := dataptr.GetAs[string]("/spec/env/"+global.EnvSystemConfig, ctx.MergedConfiguration)
	if err != nil {
		return nil, err
	}

	// does the file exist, if not, create it
	cfg, err := ReadOrCreateFile(p.fs, systemConfigPath, p.getDefault)
	if err != nil {
		return nil, err
	}

	// convert cfg to map
	m := map[string]any{}
	err = mapstructure.Decode(cfg, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *systemProvider) getDefault() Config {
	return Config{}
}
