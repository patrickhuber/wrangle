package bootstrap

import (
	"fmt"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
	"github.com/patrickhuber/wrangle/internal/install"
)

type service struct {
	install       install.Service
	configuration Configuration   `inject:"bootstrap"`
	logger        log.Logger      `inject:"logger"`
	fs            fs.FS           `inject:"fs"`
	path          filepath.Provider `inject:"filepath"`
}

type Request struct {
	Force bool
}

type Service interface {
	Execute(r *Request) error
}

func NewService(
	install install.Service,
	configuration Configuration,
	logger log.Logger,
	fs fs.FS,
	path filepath.Provider) Service {
	return &service{
		install:       install,
		configuration: configuration,
		logger:        logger,
		fs:            fs,
		path:          path,
	}
}

// Execute executes the bootstrap request
// bootstrap
//   - creates the global configuration file if it doesn't exist
//     the file is then integrated into the current configuration
//   - creates the user configuration file if it doesn't exist
//     the file is then integrated into the current configuration
//   - installs any packages in the packages directory
func (b *service) Execute(r *Request) error {

	b.logger.Debugln("bootstrap")

	cfg, err := b.configuration.Get()
	if err != nil {
		return err
	}

	// make sure the bin directory exists
	binDirectory := cfg.Spec.Environment[global.EnvBin]
	err = b.fs.MkdirAll(binDirectory, 0755)
	if err != nil {
		return fmt.Errorf("BootstrapService : failed to create bin directory %s: %w", binDirectory, err)
	}

	return b.installPackages(cfg, r.Force)
}

func (b *service) installPackages(cfg config.Config, force bool) error {
	b.logger.Debugln("install packages")
	for _, pkg := range cfg.Spec.Packages {
		request := &install.Request{
			Package: pkg.Name,
			Version: pkg.Version,
			Force:   force,
		}
		b.logger.Debugf("install %s@%s", pkg.Name, pkg.Version)
		err := b.install.Execute(request)
		if err != nil {
			return fmt.Errorf("failed to install package %s@%s: %w", pkg.Name, pkg.Version, err)
		}
	}
	return nil
}
