package config

import (
	"fmt"

	iofs "io/fs"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/global"
)

type userProvider struct {
	fs               fs.FS
	path             filepath.Provider
	errorIfNotExists bool
}

func NewUserProvider(fs fs.FS, path filepath.Provider, errorIfNotExists bool) config.Provider {
	return &userProvider{
		fs:               fs,
		path:             path,
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
			return nil, fmt.Errorf("%w user config file %s does not exist. Run `wrangle bootstrap` to create it", iofs.ErrNotExist, userConfigPath)
		}
	}

	// get the directory of the user config file
	userConfigPathDir := p.path.Dir(userConfigPath)

	// create the directory if it doesn't exist
	err = p.fs.MkdirAll(userConfigPathDir, 0755)
	if err != nil {
		return nil, err
	}

	// does the file exist, if not, create it
	cfg, err := ReadOrCreateFile(p.fs, userConfigPath, p.getDefault)
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

func (p *userProvider) getDefault() (Config, error) {
	return Config{
		ApiVersion: ApiVersion,
		Kind:       Kind,
	}, nil
}
