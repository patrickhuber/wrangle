package bootstrap

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
	fs fs.FS,
	path filepath.Provider,
	cli config.CliContext,
	systemDefaultProvider config.SystemDefaultProvider) (Configuration, error) {

	errorIfNotExists := false
	root := goconfig.NewRoot(
		config.NewEnvProvider(env, global.EnvPrefix),
		config.NewCliProvider(cli),
		config.NewSystemProvider(fs, path, systemDefaultProvider, errorIfNotExists),
		config.NewUserProvider(fs, path, errorIfNotExists),
	)
	return config.NewConfiguration(root), nil
}

func NewTestConfiguration(
	os os.OS,
	env env.Environment,
	fs fs.FS,
	path filepath.Provider,
	cli config.CliContext,
	systemDefaultProvider config.SystemDefaultProvider) Configuration {

	errorIfNotExists := false
	root := goconfig.NewRoot(
		config.NewTestDefaultProvider(os, env, path),
		config.NewEnvProvider(env, global.EnvPrefix),
		config.NewCliProvider(cli),
		config.NewSystemProvider(fs, path, systemDefaultProvider, errorIfNotExists),
		config.NewUserProvider(fs, path, errorIfNotExists),
	)
	return config.NewConfiguration(root)
}
