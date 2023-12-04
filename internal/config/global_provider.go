package config

import (
	"fmt"

	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/global"
)

// GlobalProvider Config Provider provides the global configuration
type GlobalProvider struct {
	os            os.OS
	environment   env.Environment
	path          *filepath.Processor
	globalDefault Config
	fs            fs.FS
}

func NewGlobalProvider(globalDefault Config, os os.OS, e env.Environment, path *filepath.Processor, fs fs.FS) GlobalProvider {
	return GlobalProvider{
		globalDefault: globalDefault,
		os:            os,
		environment:   e,
		path:          path,
		fs:            fs,
	}
}

func (g GlobalProvider) Get() (Config, error) {
	globalFilePath, err := g.getGlobalConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	cfg, err := NewFile(g.fs, globalFilePath).Read()
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (g GlobalProvider) Set(c Config) error {
	globalFilePath, err := g.getGlobalConfigFilePath()
	if err != nil {
		return err
	}
	return NewFile(g.fs, globalFilePath).Write(c)
}

func (g GlobalProvider) getGlobalConfigFilePath() (string, error) {
	globalConfigFilePath, ok := g.environment.Lookup(global.EnvConfig)
	if ok {
		return globalConfigFilePath, nil
	}
	globalConfigFilePath = getDefaultGlobalConfigPath(g.os, g.path)
	return globalConfigFilePath, g.environment.Set(global.EnvConfig, globalConfigFilePath)
}

func getDefaultGlobalConfigPath(os os.OS, path *filepath.Processor) string {
	return path.Join(os.Home(), ".wrangle", "config.yml")
}

func (g GlobalProvider) Default() Config {
	return g.globalDefault
}

func NewGlobalDefault(os os.OS, e env.Environment, path *filepath.Processor) (Config, error) {
	rootDirectory := "/opt/wrangle"
	plat := os.Platform()

	p := platform.Platform(plat)
	switch {
	case p.IsWindows():
		programData := e.Get("PROGRAMDATA")
		rootDirectory = path.Join(programData, "wrangle")
	case p.IsUnix():
		break
	default:
		return Config{}, fmt.Errorf("%s is unsupported", plat)
	}
	cfg := Config{
		ApiVersion: ConfigApiVersion,
		Spec: Spec{
			Environment: map[string]string{
				global.EnvBin:      path.Join(rootDirectory, "bin"),
				global.EnvConfig:   getDefaultGlobalConfigPath(os, path),
				global.EnvPackages: path.Join(rootDirectory, "packages"),
				global.EnvRoot:     rootDirectory,
			},
			Packages: []Package{
				{
					Name:    "wrangle",
					Version: "latest",
				},
				{
					Name:    "shim",
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
	}
	return cfg, nil
}

func NewGlobalTest(os os.OS, e env.Environment, path *filepath.Processor) (Config, error) {
	cfg, err := NewGlobalDefault(os, e, path)
	if err != nil {
		return Config{}, err
	}
	cfg.Spec.Feeds = []Feed{
		{
			Name: "default",
			Type: "memory",
		},
	}
	return cfg, err
}
