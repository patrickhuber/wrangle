package services

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/env"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type bootstrap struct {
	install       Install
	fs            fs.FS
	path          *filepath.Processor
	configuration Configuration
	logger        log.Logger
	environment   env.Environment
}

type BootstrapRequest struct {
	Force            bool
	RootDirectory    string
	BinDirectory     string
	ConfigFile       string
	PackageDirectory string
}

type Bootstrap interface {
	Execute(r *BootstrapRequest) error
}

func NewBootstrap(
	install Install,
	fs fs.FS,
	path *filepath.Processor,
	configuration Configuration,
	environment env.Environment,
	logger log.Logger) Bootstrap {
	return &bootstrap{
		install:       install,
		fs:            fs,
		path:          path,
		configuration: configuration,
		environment:   environment,
		logger:        logger,
	}
}

func (b *bootstrap) Execute(r *BootstrapRequest) error {
	b.logger.Debugln("bootstrap")

	// load the default configuration path
	// overwrite if the config file is set as a request parameter
	var globalConfigFilePath = r.ConfigFile
	if globalConfigFilePath == "" {
		globalConfigFilePath = b.configuration.DefaultGlobalConfigFilePath()
	}

	// fetch the global default from the configuration service we do it here so
	cfg := b.configuration.GlobalDefault()

	// overwrite any parameters specified in the request
	cfg = b.overwriteConfigDefaults(cfg, r)

	// write the changes back to the config file
	err := config.WriteFile(b.fs, globalConfigFilePath, cfg)
	if err != nil {
		return err
	}

	directories := []string{
		cfg.Spec.Environment[global.EnvPackages],
		cfg.Spec.Environment[global.EnvBin],
	}

	err = b.createDirectories(directories)
	if err != nil {
		return err
	}

	err = b.setEnvironmentVariables(cfg)
	if err != nil {
		return err
	}

	return b.installPackages(cfg)
}

func (b *bootstrap) createDirectories(directories []string) error {
	for _, dir := range directories {
		b.logger.Debugf("creating %s", dir)
		err := b.fs.MkdirAll(dir, 0775)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bootstrap) overwriteConfigDefaults(cfg config.Config, req *BootstrapRequest) config.Config {

	if req.BinDirectory != "" {
		cfg.Spec.Environment[global.EnvBin] = req.BinDirectory
	}

	if req.ConfigFile != "" {
		cfg.Spec.Environment[global.EnvConfig] = req.ConfigFile
	}

	if req.PackageDirectory != "" {
		cfg.Spec.Environment[global.EnvPackages] = req.PackageDirectory
	}

	if req.RootDirectory != "" {
		cfg.Spec.Environment[global.EnvRoot] = req.RootDirectory
	}

	return cfg
}

func (b *bootstrap) installPackages(cfg config.Config) error {
	b.logger.Debugln("install packages")
	for _, pkg := range cfg.Spec.Packages {
		request := &InstallRequest{
			Package: pkg.Name,
			Version: pkg.Version,
		}
		b.logger.Debugf("install %s@%s", pkg.Name, pkg.Version)
		err := b.install.Execute(request)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *bootstrap) setEnvironmentVariables(cfg config.Config) error {
	keys := []string{global.EnvBin, global.EnvConfig, global.EnvRoot, global.EnvPackages}
	for _, k := range keys {
		v := cfg.Spec.Environment[k]
		err := b.environment.Set(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}
