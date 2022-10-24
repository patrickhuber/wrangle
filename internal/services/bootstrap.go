package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

type bootstrap struct {
	install        Install
	initialize     Initialize
	fs             filesystem.FileSystem
	configProvider config.Provider
	logger         ilog.Logger
}

type BootstrapRequest struct {
	ApplicationName string
	Force           bool
}

type Bootstrap interface {
	Execute(r *BootstrapRequest) error
}

func NewBootstrap(
	install Install,
	initialize Initialize,
	fs filesystem.FileSystem,
	configProvider config.Provider,
	logger ilog.Logger) Bootstrap {
	return &bootstrap{
		install:        install,
		initialize:     initialize,
		fs:             fs,
		configProvider: configProvider,
		logger:         logger,
	}
}

func (b *bootstrap) Execute(r *BootstrapRequest) error {
	b.logger.Debugln("bootstrap")
	err := b.initialize.Execute(&InitializeRequest{
		ApplicationName: r.ApplicationName,
		Force:           r.Force,
	})
	if err != nil {
		return err
	}

	err = b.createDirectories()
	if err != nil {
		return err
	}
	return b.installPackages(r)
}

func (b *bootstrap) createDirectories() error {
	cfg, err := b.configProvider.Get()
	if err != nil {
		return err
	}
	directories := []string{
		cfg.Paths.Bin,
		cfg.Paths.Packages,
	}
	for _, dir := range directories {
		b.logger.Debugf("creating %s", dir)
		err = b.fs.MkdirAll(dir, 0775)
		if err != nil {
			return err
		}
	}
	return nil
}
func (b *bootstrap) installPackages(req *BootstrapRequest) error {
	b.logger.Debugln("install packages")
	references, err := b.getPackageReferences(req)
	if err != nil {
		return err
	}
	for _, reference := range references {
		request := &InstallRequest{
			Package: reference.Name,
			Version: reference.Version,
		}
		b.logger.Debugf("install %s@%s", reference.Name, reference.Version)
		err := b.install.Execute(request)
		if err != nil {
			return err
		}
	}
	return nil
}

// getPackageReferences loads the packages references from the global configuration
func (b *bootstrap) getPackageReferences(req *BootstrapRequest) ([]*config.Reference, error) {
	// maybe properties are passed into the get function?
	cfg, err := b.configProvider.Get()
	if err != nil {
		return nil, err
	}
	return cfg.References, nil
}
