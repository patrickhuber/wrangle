package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/urfave/cli/v2"
)

type Configuration interface {
	Get() (Config, error)
}

func NewConfiguration(root config.Root) Configuration {
	return &configuration{
		root: root,
	}
}

func NewDefaultConfiguration(
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	cli *cli.Context) Configuration {

	errorIfNotExists := true
	root := config.NewRoot(
		NewEnvProvider(env, global.EnvSystemConfig, global.EnvUserConfig),
		NewCliProvider(cli),
		NewSystemProvider(fs, errorIfNotExists),
		NewUserProvider(fs, errorIfNotExists),
		NewLocalProvider(fs, path, errorIfNotExists),
	)
	return &configuration{
		root: root,
	}
}

func NewBootstrapConfiguration(
	env env.Environment,
	fs fs.FS,
	cli *cli.Context) Configuration {

	errorIfNotExists := false
	root := config.NewRoot(
		NewEnvProvider(env, global.EnvSystemConfig, global.EnvUserConfig),
		NewCliProvider(cli),
		NewSystemProvider(fs, errorIfNotExists),
		NewUserProvider(fs, errorIfNotExists),
	)
	return &configuration{
		root: root,
	}
}

func NewInitializeConfiguration(
	env env.Environment,
	cli *cli.Context,
	fs fs.FS,
	path filepath.Provider,
	providers ...config.Provider) Configuration {

	root := config.NewRoot(
		NewEnvProvider(env, global.EnvSystemConfig, global.EnvUserConfig),
		NewCliProvider(cli),
		NewSystemProvider(fs, true),
		NewUserProvider(fs, true),
		NewLocalProvider(fs, path, false),
	)
	return &configuration{
		root: root,
	}
}

type configuration struct {
	root config.Root
}

func (c *configuration) Get() (Config, error) {
	data, err := c.root.Get(&config.GetContext{})
	if err != nil {
		return Config{}, err
	}
	cfg := Config{}
	err = mapstructure.Decode(data, &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}
