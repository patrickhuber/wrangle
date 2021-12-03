package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type bootstrap struct {
	install Install
	fs      filesystem.FileSystem
	cfg     *config.Config
}

type BootstrapRequest struct {
	ApplicationName  string
	GlobalConfigFile string
	Force            bool
}

type Bootstrap interface {
	Execute(r *BootstrapRequest) error
}

func NewBootstrap(i Install, fs filesystem.FileSystem, cfg *config.Config) Bootstrap {
	return &bootstrap{
		install: i,
		fs:      fs,
		cfg:     cfg,
	}
}

func (b *bootstrap) Execute(r *BootstrapRequest) error {
	err := b.validate()
	if err != nil {
		return err
	}

	err = b.createGlobalConfig(r)
	if err != nil {
		return err
	}

	return b.installPackages(r)
}

func (b *bootstrap) validate() error {
	return nil
}

func (b *bootstrap) createGlobalConfig(req *BootstrapRequest) error {
	// TODO: req.Force?
	configProvider := config.NewFileProvider(b.fs, req.GlobalConfigFile)
	return configProvider.Set(b.cfg)
}

func (b *bootstrap) installPackages(req *BootstrapRequest) error {
	packageList := []string{"wrangle", "shim"}
	for _, p := range packageList {
		request := &InstallRequest{
			Package:          p,
			GlobalConfigFile: req.GlobalConfigFile,
		}
		err := b.install.Execute(request)
		if err != nil {
			return err
		}
	}
	return nil
}
