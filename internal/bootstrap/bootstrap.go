package bootstrap

import (
	"fmt"
	"strings"

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
//   - cleans up any old renamed executables
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

	// cleanup old renamed executables
	err = b.cleanupOldExecutables(cfg)
	if err != nil {
		b.logger.Warnf("failed to cleanup old executables: %v", err)
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

func (b *service) cleanupOldExecutables(cfg config.Config) error {
	b.logger.Debugln("cleaning up old renamed executables")
	
	packagesDir := cfg.Spec.Environment[global.EnvPackages]
	
	// Walk through all package directories
	entries, err := b.fs.ReadDir(packagesDir)
	if err != nil {
		// If packages directory doesn't exist, nothing to clean up
		if strings.Contains(err.Error(), "no such file") || strings.Contains(err.Error(), "cannot find") {
			return nil
		}
		return fmt.Errorf("failed to read packages directory: %w", err)
	}

	for _, pkgEntry := range entries {
		if !pkgEntry.IsDir() {
			continue
		}
		
		pkgPath := b.path.Join(packagesDir, pkgEntry.Name())
		versionEntries, err := b.fs.ReadDir(pkgPath)
		if err != nil {
			b.logger.Warnf("failed to read package directory %s: %v", pkgPath, err)
			continue
		}

		for _, versionEntry := range versionEntries {
			if !versionEntry.IsDir() {
				continue
			}

			versionPath := b.path.Join(pkgPath, versionEntry.Name())
			fileEntries, err := b.fs.ReadDir(versionPath)
			if err != nil {
				b.logger.Warnf("failed to read version directory %s: %v", versionPath, err)
				continue
			}

			// Look for .old files and remove them
			for _, fileEntry := range fileEntries {
				if fileEntry.IsDir() {
					continue
				}
				
				if strings.HasSuffix(fileEntry.Name(), ".old") {
					oldFilePath := b.path.Join(versionPath, fileEntry.Name())
					b.logger.Debugf("removing old executable: %s", oldFilePath)
					err = b.fs.Remove(oldFilePath)
					if err != nil {
						b.logger.Warnf("failed to remove old executable %s: %v", oldFilePath, err)
					}
				}
			}
		}
	}
	
	return nil
}
