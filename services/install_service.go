package services

import (	
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

// InstallService is a service responsible for installing packages
type InstallService interface {
	Install(packagesPath string, packageName string, packageVersion string) error
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

func (service *installService) Install(packageRoot string, packageName string, packageVersion string) error {
	pkg, err := service.manager.Load(packageRoot, packageName, packageVersion)
	if err != nil{
		return err
	}
	// install the package using the manager
	return service.manager.Install(pkg)
}