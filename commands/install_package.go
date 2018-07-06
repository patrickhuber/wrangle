package commands

import (
	"fmt"
	"strings"

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
	configPackage, err := findConfigPackage(cfg, packageName)
	if err != nil {
		return err
	}
	pkg, err := cmd.createPackageFromConfig(configPackage)
	if err != nil {
		return err
	}

	// create the manager
	manager := packages.NewManager(cmd.fileSystem)

	// download the package
	fmt.Fprintf(cmd.console.Out(), "downloading '%s' to '%s'", pkg.Download().URL(), pkg.Download().OutPath())
	fmt.Fprintln(cmd.console.Out())

	err = manager.Download(pkg)
	if err != nil {
		return err
	}

	// extract the package if extraction was set
	if pkg.Extract() == nil {
		return nil
	}

	fmt.Fprintf(cmd.console.Out(), "extracting '%s' to '%s'", pkg.Download().OutPath(), pkg.Extract().OutPath())
	fmt.Fprintln(cmd.console.Out())

	return manager.Extract(pkg)
}

func findConfigPackage(cfg *config.Config, packageName string) (*config.Package, error) {
	if strings.TrimSpace(packageName) == "" {
		return nil, fmt.Errorf("package name is required")
	}
	var configPackage *config.Package
	for i := range cfg.Packages {
		pkg := &cfg.Packages[i]
		if pkg.Name == packageName {
			configPackage = pkg
			break
		}
	}
	if configPackage == nil {
		return nil, fmt.Errorf("unable to locate a package named '%s' in configuration file", packageName)
	}
	return configPackage, nil
}

func (cmd *installPackage) createPackageFromConfig(configPackage *config.Package) (packages.Package, error) {
	for i := range configPackage.Platforms {
		platform := &configPackage.Platforms[i]
		if platform.Name == cmd.platform {
			return cmd.createPackageFromPlatformConfig(configPackage, platform)
		}
	}
	return nil, fmt.Errorf("Unable to find package '%s' for platform '%s'", configPackage.Name, cmd.platform)
}

func (cmd *installPackage) createPackageFromPlatformConfig(
	configPackage *config.Package,
	platform *config.Platform) (packages.Package, error) {

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
		platform.Alias,
		download,
		extract)

	return pkg, nil
}
