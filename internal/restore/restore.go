package restore

import (
	"fmt"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/install"
)

type Request struct {
	Force bool
}

type Service interface {
	Execute(r *Request) error
}

type service struct {
	configuration config.Service
	install       install.Service
	log           log.Logger
}

func NewService(
	configuration config.Service,
	install install.Service,
	log log.Logger,
) Service {
	return &service{
		configuration: configuration,
		install:       install,
		log:           log,
	}
}

func (s *service) Execute(r *Request) error {
	cfg, err := s.configuration.Get()
	if err != nil {
		return fmt.Errorf("RestoreService : unable to get configuration: %w", err)
	}

	if len(cfg.Spec.Packages) == 0 {
		s.log.Infof("no packages found in configuration")
		return nil
	}

	var errors []error
	for _, pkg := range cfg.Spec.Packages {
		s.log.Infof("restoring package %s@%s", pkg.Name, pkg.Version)
		err := s.install.Execute(&install.Request{
			Package: pkg.Name,
			Version: pkg.Version,
			Force:   r.Force,
		})
		if err != nil {
			s.log.Warnf("failed to restore package %s@%s: %v", pkg.Name, pkg.Version, err)
			errors = append(errors, fmt.Errorf("package %s@%s: %w", pkg.Name, pkg.Version, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("RestoreService : %d package(s) failed to restore", len(errors))
	}
	return nil
}
