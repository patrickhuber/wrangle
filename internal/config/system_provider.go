package config

import (
	"fmt"
	iofs "io/fs"
	"path"

	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-dataptr"
	"github.com/patrickhuber/wrangle/internal/global"
)

func NewSystemProvider(
	fs fs.FS,
	path filepath.Provider,
	systemDefaultProvider SystemDefaultProvider,
	errorIfNotExists bool) config.Provider {
	return &systemProvider{
		fs:                    fs,
		path:                  path,
		systemDefaultProvider: systemDefaultProvider,
		errorIfNotExists:      errorIfNotExists,
	}
}

func NewTestSystemDefaultProvider(
	path filepath.Provider) SystemDefaultProvider {
	return &testSystemDefaultProvider{
		path: path,
	}
}

func NewSystemDefaultProvider(
	path filepath.Provider) SystemDefaultProvider {
	return &systemDefaultProvider{
		path: path,
	}
}

type systemProvider struct {
	fs                    fs.FS
	path                  filepath.Provider
	systemDefaultProvider SystemDefaultProvider
	errorIfNotExists      bool
}

type SystemDefaultProvider interface {
	DefaultProvider
}

type testSystemDefaultProvider struct {
	path filepath.Provider
}

type systemDefaultProvider struct {
	path filepath.Provider
}

// Get implements config.Provider.
func (p *systemProvider) Get(ctx *config.GetContext) (any, error) {
	// use a dataptr to get the value of the env config
	systemConfigPath, err := dataptr.GetAs[string]("/spec/env/"+global.EnvSystemConfig, ctx.MergedConfiguration)
	if err != nil {
		return nil, fmt.Errorf("config.SystemProvider: unable to get system config path: %w", err)
	}

	// if the system config doesn't exist and the errorIfNotExists flag is set, return an error
	if p.errorIfNotExists {
		exists, err := p.fs.Exists(systemConfigPath)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("%w system config file %s does not exist. Run `wrangle bootstrap` to create it", iofs.ErrNotExist, systemConfigPath)
		}
	}

	// get the directory of the system config file
	systemConfigDir := p.path.Dir(systemConfigPath)

	// create the directory if it doesn't exist
	err = p.fs.MkdirAll(systemConfigDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("config.SystemProvider: unable to create system config directory: %w", err)
	}

	// does the file exist, if not, create it
	cfg, err := ReadOrCreateFile(p.fs, systemConfigPath, func() (Config, error) {
		return p.systemDefaultProvider.GetDefault(ctx)
	})
	if err != nil {
		return nil, fmt.Errorf("config.SystemProvider: unable to read or create system config file: %w", err)
	}

	// convert cfg to map
	m := map[string]any{}
	err = mapstructure.Decode(cfg, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (p *testSystemDefaultProvider) GetDefault(ctx *config.GetContext) (Config, error) {
	// get the root directory from the merged configuration
	rootDirectory, err := dataptr.GetAs[string]("/spec/env/"+global.EnvRoot, ctx.MergedConfiguration)
	if err != nil {
		return Config{}, fmt.Errorf("config.SystemProvider: unable to get root directory from merged configuration: %w", err)
	}
	return Config{
		ApiVersion: ApiVersion,
		Kind:       Kind,
		Spec: Spec{
			Environment: map[string]string{
				global.EnvBin:      p.path.Join(rootDirectory, "bin"),
				global.EnvPackages: p.path.Join(rootDirectory, "packages"),
				global.EnvRoot:     rootDirectory,
			},
			Packages: []Package{
				{
					Name:    "wrangle",
					Version: "latest",
				},
			},
			Feeds: []Feed{
				{
					Name: "default",
					Type: "memory",
				},
			},
		},
	}, nil
}

func (p *systemDefaultProvider) GetDefault(ctx *config.GetContext) (Config, error) {
	// get the root directory from the merged configuration
	rootDirectory, err := dataptr.GetAs[string]("/spec/env/"+global.EnvRoot, ctx.MergedConfiguration)
	if err != nil {
		return Config{}, fmt.Errorf("config.SystemProvider: unable to get root directory from merged configuration: %w", err)
	}
	return Config{
		ApiVersion: ApiVersion,
		Kind:       Kind,
		Spec: Spec{
			Environment: map[string]string{
				global.EnvBin:      path.Join(rootDirectory, "bin"),
				global.EnvPackages: path.Join(rootDirectory, "packages"),
				global.EnvRoot:     rootDirectory,
			},
			Packages: []Package{
				{
					Name:    "wrangle",
					Version: "latest",
				},
			},
			Feeds: []Feed{
				{
					Name: "default",
					Type: "git",
					URI:  "https://github.com/patrickhuber/wrangle-packages",
				},
			},
		},
	}, nil
}
