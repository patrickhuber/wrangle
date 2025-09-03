package bootstrap

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/install"
)

type service struct {
	install       install.Service
	configuration Configuration `inject:"bootstrap"`
	logger        log.Logger
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
	logger log.Logger) Service {
	return &service{
		install:       install,
		configuration: configuration,
		logger:        logger,
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

	return b.installPackages(cfg)
}

func (b *service) installPackages(cfg config.Config) error {
	b.logger.Debugln("install packages")
	for _, pkg := range cfg.Spec.Packages {
		request := &install.Request{
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
