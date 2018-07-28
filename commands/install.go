package commands

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/patrickhuber/wrangle/filesystem"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/ui"
)

type install struct {
	platform     string
	packagesPath string
	fileSystem   filesystem.FsWrapper
	console      ui.Console
}

// Install defines an install package command
type Install interface {
	Execute(cfg *config.Config, packageName string) error
}

// NewInstall creates a new install package command
func NewInstall(
	platform string,
	packagesPath string,
	fileSystem filesystem.FsWrapper,
	console ui.Console) (Install, error) {
	if packagesPath == "" {
		return nil, fmt.Errorf("packages path can not be empty")
	}
	return &install{
			platform:     platform,
			packagesPath: packagesPath,
			fileSystem:   fileSystem,
			console:      console},
		nil
}

func (cmd *install) Execute(cfg *config.Config, packageName string) error {
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

	// initialize the link source for console output
	linkSource := pkg.Download().OutPath()

	// extract the package if extraction was set
	if pkg.Extract() != nil {
		fmt.Fprintf(cmd.console.Out(), "extracting '%s' to '%s'", pkg.Download().OutPath(), pkg.Extract().OutPath())
		fmt.Fprintln(cmd.console.Out())
		err = manager.Extract(pkg)
		if err != nil {
			return err
		}
		// the link source is now the extracted binary
		linkSource = pkg.Extract().OutPath()
	}
	linkSource = filepath.ToSlash(linkSource)

	linkTarget := filepath.Join(pkg.Download().OutFolder(), pkg.Alias())
	linkTarget = filepath.ToSlash(linkTarget)

	fmt.Fprintf(cmd.console.Out(), "linking '%s' to '%s'", linkSource, linkTarget)
	fmt.Fprintln(cmd.console.Out())

	// link the package
	return manager.Link(pkg)
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

func (cmd *install) createPackageFromConfig(configPackage *config.Package) (packages.Package, error) {
	for i := range configPackage.Platforms {
		platform := &configPackage.Platforms[i]
		if platform.Name == cmd.platform {
			return cmd.createPackageFromPlatformConfig(configPackage, platform)
		}
	}
	return nil, fmt.Errorf("Unable to find package '%s' for platform '%s'", configPackage.Name, cmd.platform)
}

func (cmd *install) createPackageFromPlatformConfig(
	configPackage *config.Package,
	platform *config.Platform) (packages.Package, error) {

	if platform.Download == nil {
		return nil, fmt.Errorf("platform 'download' element is required for package '%s' platform '%s'", configPackage.Name, platform.Name)
	}

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
