package commands

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/filesystem"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/patrickhuber/cli-mgr/packages"
	"github.com/patrickhuber/cli-mgr/ui"
)

type installPackage struct {
	platform     string
	packagesPath string
	fileSystem   filesystem.FsWrapper
	console      ui.Console
}

// InstallPackage defines an install package command
type InstallPackage interface {
	Execute(cfg *config.Config, packageName string) error
}

// NewInstallPackage creates a new install package command
func NewInstallPackage(
	platform string,
	packagesPath string,
	fileSystem filesystem.FsWrapper,
	console ui.Console) InstallPackage {
	return &installPackage{
		platform:     platform,
		packagesPath: packagesPath,
		fileSystem:   fileSystem,
		console:      console}
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
				cmd.packagesPath,
				platform.Download.Out)

			// create the extraction part of the package
			var extract packages.Extract
			if platform.Extract != nil {
				extract = packages.NewExtract(
					platform.Extract.Filter,
					cmd.packagesPath,
					platform.Extract.Out)
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
			fmt.Fprintf(cmd.console.Out(), "downloading '%s' to '%s'", pkg.Download().URL(), pkg.Download().OutPath())
			fmt.Fprintln(cmd.console.Out())

			err := manager.Download(pkg)
			if err != nil {
				return err
			}

			// extract the package if extraction was set
			if extract != nil {

				fmt.Fprintf(cmd.console.Out(), "extracting '%s' to '%s'", pkg.Download().OutPath(), pkg.Extract().OutPath())
				fmt.Fprintln(cmd.console.Out())

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
