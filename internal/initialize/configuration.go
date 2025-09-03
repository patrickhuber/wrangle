package initialize

import (
	goconfig "github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Configuration interface {
	config.Configuration
}

func NewConfiguration(
	env env.Environment,
	cli config.CliContext,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	globResolver goconfig.GlobResolver,
	systemDefaultProvider config.SystemDefaultProvider) (Configuration, error) {

	builder := goconfig.NewBuilder()
	builder.WithProvider(config.NewEnvProvider(env, global.EnvPrefix))
	builder.WithProvider(config.NewCliProvider(cli))
	builder.WithProvider(config.NewSystemProvider(fs, path, systemDefaultProvider, true))
	builder.WithProvider(config.NewUserProvider(fs, path, true))
	builder.WithFactory(config.NewLocalFactory(fs, path, os, globResolver, false))

	root, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return config.NewConfiguration(root), nil
}

func NewTestConfiguration(
	env env.Environment,
	cli config.CliContext,
	fs fs.FS,
	path filepath.Provider,
	os os.OS,
	globResolver goconfig.GlobResolver,
	systemDefaultProvider config.SystemDefaultProvider) (Configuration, error) {

	builder := goconfig.NewBuilder()
	builder.WithProvider(config.NewTestDefaultProvider(os, env, path))
	builder.WithProvider(config.NewEnvProvider(env, global.EnvPrefix))
	builder.WithProvider(config.NewCliProvider(cli))
	builder.WithProvider(config.NewSystemProvider(fs, path, systemDefaultProvider, true))
	builder.WithProvider(config.NewUserProvider(fs, path, true))
	builder.WithFactory(config.NewLocalFactory(fs, path, os, globResolver, false))

	root, err := builder.Build()
	if err != nil {
		return nil, err
	}
	return config.NewConfiguration(root), nil
}
