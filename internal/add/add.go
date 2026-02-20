package add

import (
	"fmt"

	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type Request struct {
	Package string
	Version string
}

type Service interface {
	Execute(r *Request) error
}

type service struct {
	fs    fs.FS
	opsys os.OS
	path  filepath.Provider
	log   log.Logger
}

func NewService(
	fs fs.FS,
	opsys os.OS,
	path filepath.Provider,
	log log.Logger,
) Service {
	return &service{
		fs:    fs,
		opsys: opsys,
		path:  path,
		log:   log,
	}
}

func (s *service) Execute(r *Request) error {
	if r.Package == "" {
		return fmt.Errorf("package name is required")
	}

	workDir, err := s.opsys.WorkingDirectory()
	if err != nil {
		return fmt.Errorf("unable to get working directory: %w", err)
	}

	localFile := s.path.Join(workDir, global.LocalConfigurationFileName)

	exists, err := s.fs.Exists(localFile)
	if err != nil {
		return fmt.Errorf("unable to check for local config file: %w", err)
	}

	var cfg config.Config
	if exists {
		cfg, err = config.ReadFile(s.fs, localFile)
		if err != nil {
			return fmt.Errorf("unable to read local config file '%s': %w", localFile, err)
		}
	} else {
		s.log.Infof("creating local config file '%s'", localFile)
		cfg = config.Config{
			ApiVersion: config.ApiVersion,
			Kind:       config.Kind,
			Spec: config.Spec{
				Feeds:       []config.Feed{},
				Stores:      []config.Store{},
				Environment: map[string]string{},
				Packages:    []config.Package{},
			},
		}
	}

	// update or add the package entry
	updated := false
	for idx, pkg := range cfg.Spec.Packages {
		if pkg.Name == r.Package {
			cfg.Spec.Packages[idx].Version = r.Version
			updated = true
			break
		}
	}
	if !updated {
		cfg.Spec.Packages = append(cfg.Spec.Packages, config.Package{
			Name:    r.Package,
			Version: r.Version,
		})
	}

	err = config.WriteFile(s.fs, localFile, cfg)
	if err != nil {
		return fmt.Errorf("unable to write local config file '%s': %w", localFile, err)
	}

	if updated {
		s.log.Infof("updated package %s@%s in '%s'", r.Package, r.Version, localFile)
	} else {
		s.log.Infof("added package %s@%s to '%s'", r.Package, r.Version, localFile)
	}
	return nil
}
