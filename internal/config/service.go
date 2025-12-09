package config

import (
	"github.com/go-viper/mapstructure/v2"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Service interface {
	Get() (Config, error)
}

func NewTestService(
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	cli CliContext,
	globResolver config.GlobResolver,
	systemDefaultProvider SystemDefaultProvider) (Service, error) {

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

func NewDefaultService(
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	cli CliContext,
	globResolver config.GlobResolver,
	systemDefaultProvider SystemDefaultProvider) (Service, error) {

	errorIfNotExists := true
	builder := config.NewBuilder()
	builder.WithProvider(NewDefaultProvider(os, env, path))
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

func NewConfiguration(root config.Root) Service {
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
