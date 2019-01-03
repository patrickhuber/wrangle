package services

import (	
	"github.com/patrickhuber/wrangle/filepath"
	"strings"
	"fmt"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

type InstallServiceRequestPackage struct{
	Name string
	Version string
}

type InstallServiceRequestDirectories struct{
	Root string
	Bin string
	Packages string
}

type InstallServiceRequestFeed struct{
	URL string
}

// InstallServiceRequest defines a request for installation of a package
type InstallServiceRequest struct{
	Package *InstallServiceRequestPackage
	Directories *InstallServiceRequestDirectories
	Feed *InstallServiceRequestFeed
}

// InstallService is a service responsible for installing packages
type InstallService interface {
	Install(request *InstallServiceRequest) error
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

func (service *installService) Install(request *InstallServiceRequest) error {	
	packageName, err := service.getPackageName(request)
	if err != nil{
		return err
	}

	packageVersion := service.getPackageVersion(request)

	root, rootIsMissing := service.getRoot(request)

	bin, err := service.getBin(request, rootIsMissing)	
	if err != nil{
		return err
	}

	packagesRoot, err := service.getPackagesRoot(request, rootIsMissing)	
	if err != nil{
		return err
	}

	pkg, err := service.manager.Load(root, bin, packagesRoot, packageName, packageVersion)
	if err != nil{
		return err
	}

	// install the package using the manager
	return service.manager.Install(pkg)
}

func  (service *installService) getPackageName(request *InstallServiceRequest) (string, error){
	packageNameIsMissing := request.Package == nil || strings.TrimSpace(request.Package.Name) == ""
	if packageNameIsMissing {
		return "", fmt.Errorf("missing required argument package name")
	}
	return request.Package.Name, nil
}

func (service *installService) getPackageVersion(request *InstallServiceRequest) string{	
	if request.Package != nil{
		return strings.TrimSpace(request.Package.Version)
	}
	return ""
}

func (service *installService) getRoot(request *InstallServiceRequest) (string, bool){	
	if request.Directories == nil || strings.TrimSpace(request.Directories.Root) == ""	{
		return "", true		
	}
	return request.Directories.Root, false
}

func (service *installService) getBin(request *InstallServiceRequest, rootIsMissing bool)(string, error){
	binIsMissing := request.Directories == nil ||  strings.TrimSpace(request.Directories.Bin) == ""	
	if !binIsMissing{
		return request.Directories.Bin, nil
	}
	if rootIsMissing {
		return "", fmt.Errorf("bin must be specified if root is not specified")
	}
	return filepath.Join(request.Directories.Root, "bin"), nil
}

func (service *installService) getPackagesRoot(request *InstallServiceRequest, rootIsMissing bool)(string, error){
	packagesRootIsMissing := request.Directories== nil || strings.TrimSpace(request.Directories.Packages) == ""
	if !packagesRootIsMissing{
		return request.Directories.Packages,nil
	}
	if rootIsMissing {
		return "", fmt.Errorf("packages root must be specified if root is not specified")
	}
	return filepath.Join(request.Directories.Root, "packages"), nil
}