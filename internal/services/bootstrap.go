package services

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type bootstrap struct {
	install Install
	fs      filesystem.FileSystem
	reader  config.Reader
}

type BootstrapRequest struct {
	ApplicationName  string
	GlobalConfigFile string
	Force            bool
}

type Bootstrap interface {
	Execute(r *BootstrapRequest) error
}

func NewBootstrap(i Install, fs filesystem.FileSystem, defaultReader config.Reader) Bootstrap {
	return &bootstrap{
		install: i,
		fs:      fs,
		reader:  defaultReader,
	}
}

func (b *bootstrap) Execute(r *BootstrapRequest) error {

	err := b.ensureGlobalConfig(r)
	if err != nil {
		return err
	}

	return b.installPackages(r)
}

func (b *bootstrap) ensureGlobalConfig(req *BootstrapRequest) error {

	// does the global config exist?
	exists, err := b.fs.Exists(req.GlobalConfigFile)
	if err != nil {
		return err
	}

	if exists && !req.Force {
		return nil
	}

	cfg, err := b.reader.Get()
	if err != nil {
		return err
	}

	configProvider := config.NewFileProvider(b.fs, req.GlobalConfigFile)
	return configProvider.Set(cfg)
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
