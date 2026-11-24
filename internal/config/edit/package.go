package edit

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/config"
)

type PackageRequest struct {
	Package string
	Version string
	Scope   string
}

func (s *service) AddPackage(r *PackageRequest) error {
	if r.Package == "" {
		return fmt.Errorf("package name is required")
	}
	version := r.Version
	if version == "" {
		version = config.TagLatest
	}

	updated := false
	file, err := s.update(r.Scope, func(cfg *config.Config) error {
		for idx, pkg := range cfg.Spec.Packages {
			if pkg.Name == r.Package {
				cfg.Spec.Packages[idx].Version = version
				updated = true
				break
			}
		}
		if !updated {
			cfg.Spec.Packages = append(cfg.Spec.Packages, config.Package{
				Name:    r.Package,
				Version: version,
			})
		}
		return nil
	})
	if err != nil {
		return err
	}
	if updated {
		s.log.Infof("updated package %s@%s in '%s'", r.Package, version, file)
	} else {
		s.log.Infof("added package %s@%s to '%s'", r.Package, version, file)
	}
	return nil
}
