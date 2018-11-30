package services

import (
	"github.com/patrickhuber/wrangle/filepath"
	"strings"
	"github.com/spf13/afero"
	"github.com/patrickhuber/wrangle/tasks"
	"fmt"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
)

// InstallService is a service responsible for installing packages
type InstallService interface {
	Install(packagesPath string, packageName string, packageVersion string) error
}

type installService struct {	
	platform     string
	fileSystem   filesystem.FsWrapper
	manager      packages.Manager
	loader       config.Loader
}

// NewInstallService creates a new install service
func NewInstallService(
	platform string,
	fileSystem filesystem.FsWrapper,
	manager packages.Manager,
	loader config.Loader) (InstallService, error) {

	return &installService{
		platform:     platform,
		fileSystem:   fileSystem,
		manager:      manager,
		loader:       loader},	nil
}

func (service *installService) Install(packagesPath string, packageName string, packageVersion string) error {

	// load the configuration for the package into the struct
	configPkg, err := service.findConfigPackage(packagesPath, packageName, packageVersion)
	if err != nil {
		return err
	}

	// turn the config into package object
	pkg, err := service.createPackageFromConfig(configPkg, service.platform)
	if err != nil {
		return err
	}
	return service.manager.Install(pkg)
}


func (service *installService) findConfigPackage(packagesPath, packageName, packageVersion string) (*config.Package, error) {
	if strings.TrimSpace(packageName) == "" {
		return nil, fmt.Errorf("package name is required")
	}

	packagePath := filepath.Join(packagesPath, packageName)
	useLatestVersion := len(strings.TrimSpace(packageVersion)) == 0

	var packageManifestPath string
	if !useLatestVersion {
		packageManifestFileName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)
		packageManifestPath = filepath.Join(packagePath, packageVersion, packageManifestFileName)

	} else {
		files, err := afero.ReadDir(service.fileSystem, packagePath)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if !file.IsDir() {
				continue
			}
			packageVersion = file.Name()
			break
		}
	}

	configPackage, err := service.loader.LoadPackage(packageManifestPath)
	if err != nil {
		return nil, err
	}
	return configPackage, nil
}

func (service *installService) createPackageFromConfig(
	configPackage *config.Package, platformName string) (packages.Package, error) {

	taskList := make([]tasks.Task, 0)

	for _, configPlatform := range configPackage.Platforms {
		if configPlatform.Name != platformName {
			continue
		}
		for _, configTask := range configPlatform.Tasks {
			params := map[string]string{}
			for key, value := range configTask.Params {
				params[key] = value.(string)
			}
			task := tasks.NewTask(configTask.Name, configTask.Type, params)
			taskList = append(taskList, task)
		}
	}
	// create the package with the download and extract params set
	pkg := packages.New(
		configPackage.Name,
		configPackage.Version,
		taskList...)

	return pkg, nil
}
