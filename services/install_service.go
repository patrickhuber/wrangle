package services

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/packages"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/templates"
	"github.com/spf13/afero"
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

	// load the configuration for the package into the struct
	configPkg, err := service.findConfigPackage(packageRoot, packageName, packageVersion)
	if err != nil {
		return err
	}

	// turn the config into package object
	pkg, err := service.createPackageFromConfig(configPkg, service.platform)
	if err != nil {
		return err
	}

	// install the package using the manager
	return service.manager.Install(pkg)
}

func (service *installService) findConfigPackage(packageRoot, packageName, packageVersion string) (*config.Package, error) {

	// packageRoot
	//   ex: /packages
	// packagePath
	//   ex: /packages/test
	// packageVersionPath
	//   ex: /packages/test/1.0.0/
	// packageVersionManifestPath
	//   ex: /packages/test/1.0.0/test.1.0.0.yml
	packagePath, err := service.findPackagePath(packageRoot, packageName)
	if err != nil {
		return nil, err
	}

	packageVersion, err = service.findPackageVersion(packagePath, packageVersion)
	if err != nil {
		return nil, err
	}

	packageVersionPath := fmt.Sprintf("%s/%s", packagePath, packageVersion)
	packageVersionManifestPath := fmt.Sprintf("%s/%s.%s.yml", packageVersionPath, packageName, packageVersion)

	packageData, err := service.loader.LoadPackageAsInterface(packageVersionManifestPath)
	if err != nil {
		return nil, err
	}

	template := templates.NewTemplate(packageData)
	dictionary := collections.NewDictionary()
	dictionary.Set("/package_install_directory", packageVersionPath)
	dictionary.Set("package_cache_directory", packagePath)

	packageData, err = template.Evaluate(templates.NewDictionaryResolver(dictionary))
	if err != nil {
		return nil, err
	}

	packageDataString, err := config.SerializePackageFromInterface(packageData)
	if err != nil {
		return nil, err
	}

	return config.DeserializePackageString(packageDataString)
}

func (service *installService) findPackagePath(packageRoot, packageName string) (string, error) {
	if strings.TrimSpace(packageName) == "" {
		return "", fmt.Errorf("package name is required")
	}

	packagePath := filepath.Join(packageRoot, packageName)
	return packagePath, nil
}

func (service *installService) findPackageVersion(packagePath, packageVersion string) (string, error) {
	useLatestVersion := len(strings.TrimSpace(packageVersion)) == 0

	if !useLatestVersion {
		return packageVersion, nil
	}

	return service.findLatestPackageVersion(packagePath)
}

func (service *installService) findLatestPackageVersion(packagePath string) (string, error) {
	files, err := afero.ReadDir(service.fileSystem, packagePath)
	if err != nil {
		return "", err
	}
	packageVersion := ""
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		packageVersion = file.Name()
		break
	}
	return packageVersion, err
}

func (service *installService) createPackageFromConfig(
	configPackage *config.Package, platformName string) (packages.Package, error) {

	taskList := make([]tasks.Task, 0)

	for _, target := range configPackage.Targets {
		if target.Platform != platformName {
			continue
		}
		/* for _, configTask := range target.Tasks {
			params := map[string]string{}
			var task tasks.Task = nil
			taskList = append(taskList, task)
		} */
	}
	// create the package with the download and extract params set
	pkg := packages.New(
		configPackage.Name,
		configPackage.Version,
		taskList...)

	return pkg, nil
}
