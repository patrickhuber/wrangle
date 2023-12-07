package services

import (
	"fmt"
	"regexp"

	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/os"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Configuration struct {
	path          *filepath.Processor
	os            os.OS
	e             env.Environment
	fs            fs.FS
	globalDefault config.Config
}

type GlobalDefault interface {
	Get() (config.Config, error)
}

func NewConfiguration(os os.OS, e env.Environment, fs fs.FS, path *filepath.Processor) (Configuration, error) {
	globalDefault, err := NewGlobalDefault(os, e, path)
	if err != nil {
		return Configuration{}, err
	}
	return Configuration{
		os:            os,
		path:          path,
		e:             e,
		fs:            fs,
		globalDefault: globalDefault,
	}, nil
}

func NewTestConfiguration(os os.OS, e env.Environment, fs fs.FS, path *filepath.Processor) (Configuration, error) {
	globalDefault, err := NewGlobalTestDefault(os, e, path)
	if err != nil {
		return Configuration{}, err
	}
	return Configuration{
		os:            os,
		path:          path,
		e:             e,
		fs:            fs,
		globalDefault: globalDefault,
	}, nil
}

func (c Configuration) Get() (config.Config, error) {
	globalConfig, err := c.GlobalConfiguration()
	if err != nil {
		return config.Config{}, err
	}

	localConfigurations, err := c.LocalConfigurations()
	if err != nil {
		return config.Config{}, err
	}

	return c.merge(globalConfig, localConfigurations...)
}

func (c Configuration) merge(global config.Config, locals ...config.Config) (config.Config, error) {
	if len(locals) == 0 {
		return global, nil
	}
	current := global

	feeds := map[string]config.Feed{}
	for _, f := range current.Spec.Feeds {
		feeds[f.Name] = f
	}

	// TODO: add a store name to avoid merging on type
	stores := map[string]config.Store{}
	for _, s := range current.Spec.Stores {
		stores[s.Type] = s
	}

	packages := map[string]config.Package{}
	for _, p := range current.Spec.Packages {
		key := p.Name
		if p.Version != "" {
			key = fmt.Sprintf("%s@%s", p.Name, p.Version)
		}
		packages[key] = p
	}

	for _, local := range locals {
		if current.ApiVersion != local.ApiVersion {
			return config.Config{}, fmt.Errorf("unable to merge configurations incompatible api versions '%s' and '%s'", current.ApiVersion, local.ApiVersion)
		}

		// apply local metadata to global metadata, overwriting any duplicates
		for k, v := range local.Metadata {
			current.Metadata[k] = v
		}

		// apply local environment to current metadata, overwriting any duplicates
		for k, v := range local.Spec.Environment {
			current.Spec.Environment[k] = v
		}

		for _, f := range local.Spec.Feeds {
			feeds[f.Name] = f
		}

		for _, s := range local.Spec.Stores {
			stores[s.Type] = s
		}

		for _, p := range local.Spec.Packages {
			key := p.Name
			if p.Version != "" {
				key = fmt.Sprintf("%s@%s", p.Name, p.Version)
			}
			packages[key] = p
		}
	}

	if len(current.Metadata) == 0 {
		current.Metadata = nil
	}

	if len(current.Spec.Environment) == 0 {
		current.Spec.Environment = nil
	}

	clear(current.Spec.Feeds)
	if len(feeds) == 0 {
		current.Spec.Feeds = nil
	}
	for _, f := range feeds {
		current.Spec.Feeds = append(current.Spec.Feeds, f)
	}

	clear(current.Spec.Stores)
	if len(stores) == 0 {
		current.Spec.Stores = nil
	}
	for _, s := range stores {
		current.Spec.Stores = append(current.Spec.Stores, s)
	}

	clear(current.Spec.Packages)
	if len(packages) == 0 {
		current.Spec.Packages = nil
	}
	for _, p := range packages {
		current.Spec.Packages = append(current.Spec.Packages, p)
	}

	return current, nil
}

func (c Configuration) GlobalDefault() config.Config {
	return c.globalDefault
}

func (c Configuration) DefaultGlobalConfigFilePath() string {
	return c.path.Join(c.os.Home(), ".wrangle", "config.yml")
}

func (c Configuration) DefaultLocalConfigFilePath() (string, error) {
	wd, err := c.os.WorkingDirectory()
	if err != nil {
		return "", err
	}
	return c.path.Join(wd, ".wrangle.yml"), nil
}

func NewGlobalDefault(os os.OS, e env.Environment, path *filepath.Processor) (config.Config, error) {
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
		return config.Config{}, fmt.Errorf("%s is unsupported", plat)
	}
	cfg := config.Config{
		ApiVersion: config.ConfigApiVersion,
		Spec: config.Spec{
			Environment: map[string]string{
				global.EnvBin:      path.Join(rootDirectory, "bin"),
				global.EnvPackages: path.Join(rootDirectory, "packages"),
				global.EnvRoot:     rootDirectory,
			},
			Packages: []config.Package{
				{
					Name:    "wrangle",
					Version: "latest",
				},
				{
					Name:    "shim",
					Version: "latest",
				},
			},
			Feeds: []config.Feed{
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

func NewGlobalTestDefault(os os.OS, e env.Environment, path *filepath.Processor) (config.Config, error) {
	cfg, err := NewGlobalDefault(os, e, path)
	if err != nil {
		return config.Config{}, err
	}
	cfg.Spec.Feeds = []config.Feed{
		{
			Name: "default",
			Type: "memory",
		},
	}
	return cfg, nil
}

const (
	LocalConfigFilePattern = "[.]wrangle[.](yml|yaml|json)"
)

func (c Configuration) LocalConfigurations() ([]config.Config, error) {

	// look in the current directory
	pwd, err := c.os.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	// work up the directory hierarchy to find every path to the root
	current := pwd
	dirs := []string{}
	for {
		dirs = append(dirs, current)

		parent := c.path.Dir(current)
		if parent == current {
			break
		}

		current = parent
	}

	// loop through all the directories looking for configuration files
	var cfgs []config.Config
	for _, dir := range dirs {

		files, err := c.fs.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		reg := regexp.MustCompile(LocalConfigFilePattern)

		// for each path match only files that match the config file pattern
		for _, file := range files {

			if file.IsDir() {
				continue
			}

			if !reg.MatchString(file.Name()) {
				continue
			}

			filePath := c.path.Join(dir, file.Name())

			// load the configuration file
			cfg, err := config.ReadFile(c.fs, filePath)
			if err != nil {
				return nil, err
			}

			cfgs = append(cfgs, cfg)
		}
	}
	return cfgs, nil
}

func (c Configuration) GlobalConfiguration() (config.Config, error) {

	// check if the env var is set, use the value if so
	globalDefault, ok := c.e.Lookup(global.EnvConfig)
	if !ok {
		// otherwise use default
		globalDefault = c.DefaultGlobalConfigFilePath()
	}

	// load the file
	cfg, err := config.ReadFile(c.fs, globalDefault)
	if err != nil {
		return config.Config{}, fmt.Errorf("%w : unable to load global configuration file '%s'", err, globalDefault)
	}
	return cfg, nil
}
