package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type bootstrap struct {
	install    Install
	initialize Initialize
	fs         filesystem.FileSystem
}

type BootstrapRequest struct {
	ApplicationName  string
	GlobalConfigFile string
	Force            bool
}

type Bootstrap interface {
	Execute(r *BootstrapRequest) error
}

func NewBootstrap(install Install, initialize Initialize, fs filesystem.FileSystem) Bootstrap {
	return &bootstrap{
		install:    install,
		initialize: initialize,
		fs:         fs,
	}
}

func (b *bootstrap) Execute(r *BootstrapRequest) error {

	if err := b.initialize.Execute(&InitializeRequest{
		ApplicationName:  r.ApplicationName,
		GlobalConfigFile: r.GlobalConfigFile,
		Force:            r.Force,
	}); err != nil {
		return err
	}

	return b.installPackages(r)
}

func (b *bootstrap) installPackages(req *BootstrapRequest) error {

	references, err := b.getPackageReferences(req)
	if err != nil {
		return err
	}
	for _, reference := range references {
		request := &InstallRequest{
			Package:          reference.Name,
			Version:          reference.Version,
			GlobalConfigFile: req.GlobalConfigFile,
		}
		err := b.install.Execute(request)
		if err != nil {
			return err
		}
	}
	return nil
}

// getPackageReferences loads the packages references from the global configuration
func (b *bootstrap) getPackageReferences(req *BootstrapRequest) ([]*config.Reference, error) {
	configProvider := config.NewFileProvider(b.fs, req.GlobalConfigFile)
	cfg, err := configProvider.Get()
	if err != nil {
		return nil, err
	}
	return cfg.References, nil
}
