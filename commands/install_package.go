package commands

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/packages"
	"github.com/spf13/afero"
)

type installPackage struct {
	platform     string
	packagesPath string
	fileSystem   afero.Fs
}

// InstallPackage defines an install package command
type InstallPackage interface {
	Execute(cfg *config.Config, packageName string) error
}

// NewInstallPackage creates a new install package command
func NewInstallPackage(
	platform string,
	packagesPath string,
	fileSystem afero.Fs) InstallPackage {
	return &installPackage{
		platform:     platform,
		packagesPath: packagesPath,
		fileSystem:   fileSystem}
}

func (cmd *installPackage) Execute(cfg *config.Config, packageName string) error {
	var configPackage *config.Package
	for i := range cfg.Packages {
		pkg := &cfg.Packages[i]
		if pkg.Name == packageName {
			configPackage = pkg
			break
		}
	}
	if configPackage == nil {
		return fmt.Errorf("unable to locate a package named '%s' in configuration file", packageName)
	}
	for _, platform := range configPackage.Platforms {
		if platform.Name == cmd.platform {

			// create the download part of the package
			download := packages.NewDownload(
				platform.Download.URL,
				platform.Download.Out,
				cmd.packagesPath)

			// create the extraction part of the package
			var extract packages.Extract
			if platform.Extract != nil {
				extract = packages.NewExtract(
					platform.Extract.Filter,
					platform.Extract.Out,
					cmd.packagesPath)
			}

			// create the package with the download and extract params set
			pkg := packages.New(
				configPackage.Name,
				configPackage.Version,
				configPackage.Alias,
				download,
				extract)

			// create the manager
			manager := packages.NewManager(cmd.fileSystem)

			// download the package
			err := manager.Download(pkg)
			if err != nil {
				return err
			}

			// extract the package if extraction was set
			if extract != nil {
				err = manager.Extract(pkg)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}
	return fmt.Errorf("unable to find platform '%s' in configuration for package '%s'", cmd.platform, configPackage.Name)
}
