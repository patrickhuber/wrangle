package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/wrangle/internal/global"
)

type DefaultProvider interface {
	GetDefault(ctx *config.GetContext) (Config, error)
}

func NewDefaultProvider(os os.OS, env env.Environment, path filepath.Provider) config.Provider {
	return &defaultProvider{
		os:   os,
		env:  env,
		path: path,
	}
}

type defaultProvider struct {
	path filepath.Provider
	os   os.OS
	env  env.Environment
}

func (p *defaultProvider) Get(ctx *config.GetContext) (any, error) {
	root, err := GetRoot(p.env, p.os.Platform())
	if err != nil {
		return nil, err
	}

	home, err := p.os.Home()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"spec": map[string]any{
			"env": map[string]any{
				global.EnvRoot:         root,
				global.EnvBin:          GetDefaultBinPath(p.path, root),
				global.EnvPackages:     GetDefaultPackagesPath(p.path, root),
				global.EnvSystemConfig: GetDefaultSystemConfigPath(p.path, root),
				global.EnvUserConfig:   GetDefaultUserConfigPath(p.path, home),
			},
		},
	}, nil
}
