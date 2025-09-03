package config

import (
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Configuration interface {
	Get() (Config, error)
}

func NewTestConfiguration(
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	cli CliContext,
	globResolver config.GlobResolver,
	systemDefaultProvider SystemDefaultProvider) (Configuration, error) {

	errorIfNotExists := true

	// the CliContext will provide some defaults in the default case that we need to supply here
	builder := config.NewBuilder()
	builder.WithProvider(NewTestDefaultProvider(os, env, path))
	builder.WithProvider(NewEnvProvider(env, global.EnvPrefix))
	builder.WithProvider(NewCliProvider(cli))
	builder.WithProvider(NewSystemProvider(fs, path, systemDefaultProvider, errorIfNotExists))
	builder.WithProvider(NewUserProvider(fs, path, errorIfNotExists))
	builder.WithFactory(NewLocalFactory(fs, path, os, globResolver, errorIfNotExists))

	root, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return NewConfiguration(root), nil
}

func NewDefaultConfiguration(
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	cli CliContext,
	globResolver config.GlobResolver,
	systemDefaultProvider SystemDefaultProvider) (Configuration, error) {

	errorIfNotExists := true
	builder := config.NewBuilder()
	builder.WithProvider(NewEnvProvider(env, global.EnvPrefix))
	builder.WithProvider(NewCliProvider(cli))
	builder.WithProvider(NewSystemProvider(fs, path, systemDefaultProvider, errorIfNotExists))
	builder.WithProvider(NewUserProvider(fs, path, errorIfNotExists))
	builder.WithFactory(NewLocalFactory(fs, path, os, globResolver, errorIfNotExists))

	root, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return &configuration{
		root: root,
	}, nil
}

type configuration struct {
	root config.Root
}

func NewConfiguration(root config.Root) Configuration {
	return &configuration{
		root: root,
	}
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
