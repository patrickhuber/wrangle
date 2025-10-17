package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/wrangle/internal/global"
)

func NewTestDefaultProvider(os os.OS, env env.Environment, path filepath.Provider) config.Provider {
	return &testDefaultProvider{
		os:   os,
		env:  env,
		path: path,
	}
}

type testDefaultProvider struct {
	path filepath.Provider
	os   os.OS
	env  env.Environment
}

// Get implements config.Provider.
func (t *testDefaultProvider) Get(ctx *config.GetContext) (any, error) {

	root, err := GetRoot(t.env, t.path, t.os.Platform())
	if err != nil {
		return nil, err
	}

	home, err := t.os.Home()
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"spec": map[string]any{
			"env": map[string]any{
				global.EnvRoot:         root,
				global.EnvBin:          GetDefaultBinPath(t.path, root),
				global.EnvPackages:     GetDefaultPackagesPath(t.path, root),
				global.EnvSystemConfig: GetDefaultSystemConfigPath(t.path, root),
				global.EnvUserConfig:   GetDefaultUserConfigPath(t.path, home),
			},
		},
	}, nil
}
