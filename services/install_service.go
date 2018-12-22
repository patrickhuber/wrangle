package services

import (	
	"github.com/patrickhuber/wrangle/filepath"
	"strings"
	"fmt"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

// InstallService is a service responsible for installing packages
type InstallService interface {
	Install(root, bin, packagesRoot , packageName , packageVersion string) error
}

type installService struct {
	platform   string
	fileSystem filesystem.FsWrapper
	manager    packages.Manager
	loader     config.Loader
}

// NewInstallService creates a new install service
func NewInstallService(
	platform string,
	fileSystem filesystem.FsWrapper,
	manager packages.Manager,
	loader config.Loader) (InstallService, error) {

	return &installService{
		platform:   platform,
		fileSystem: fileSystem,
		manager:    manager,
		loader:     loader}, nil
}

func (service *installService) Install(root, bin, packagesRoot, packageName, packageVersion string) error {	

	packageNameIsMissing := strings.TrimSpace(packageName) == ""
	if packageNameIsMissing {
		return fmt.Errorf("missing required argument package name")
	}

	rootIsMissing := strings.TrimSpace(root) == ""
	binIsMissing := strings.TrimSpace(bin) == ""
	packagesRootIsMissing := strings.TrimSpace(packagesRoot) == ""

	if rootIsMissing && binIsMissing {
		return fmt.Errorf("bin must be specified if root is not specified")
	}else if binIsMissing{
		bin = filepath.Join(root, "bin")
	}

	if rootIsMissing && packagesRootIsMissing {
		return fmt.Errorf("packages root must be specified if root is not specified")
	}else if packagesRootIsMissing{
		packagesRoot = filepath.Join(root, "packages")
	}

	pkg, err := service.manager.Load(root, bin, packagesRoot, packageName, packageVersion)
	if err != nil{
		return err
	}
	// install the package using the manager
	return service.manager.Install(pkg)
}