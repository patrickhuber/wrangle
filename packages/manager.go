package packages

import (	
	semver "github.com/hashicorp/go-version"
	"github.com/patrickhuber/wrangle/filepath"
	"fmt"
	"github.com/spf13/afero"
	"strings"
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/templates"
)

type manager struct {
	fileSystem    filesystem.FsWrapper
	taskProviders tasks.ProviderRegistry
}

// Manager defines a manager interface
type Manager interface {
	Install(p Package) error
	Load(root, bin, packagesRoot, packageName, packageVersion string) (Package, error)
}

// NewManager creates a new package manager
func NewManager(fileSystem filesystem.FsWrapper, taskProviders tasks.ProviderRegistry) Manager {
	return &manager{
		fileSystem:    fileSystem,
		taskProviders: taskProviders}
}

func (manager *manager) Install(p Package) error {
	for _, task := range p.Tasks() {
		provider, err := manager.taskProviders.Get(task.Type())
		if err != nil {
			return err
		}
		err = provider.Execute(task, p.Context())
		if err != nil {
			return err
		}
	}
	return nil
}

func (manager *manager) Load(root, bin, packagesRoot, packageName, packageVersion string) (Package, error) {

	// packagesRoot
	//   ex: /packages
	// packagePath
	//   ex: /packages/test
	// packageVersionPath
	//   ex: /packages/test/1.0.0/
	// packageVersionManifestPath
	//   ex: /packages/test/1.0.0/test.1.0.0.yml
	packagePath, err := manager.getPackagePath(packagesRoot, packageName)
	if err != nil {
		return nil, err
	}

	packageVersion, err = manager.getPackageVersion(packagePath, packageVersion)
	if err != nil {
		return nil, err
	}

	packageVersionPath := fmt.Sprintf("%s/%s", packagePath, packageVersion)
	packageVersionManifestPath := fmt.Sprintf("%s/%s.%s.yml", packageVersionPath, packageName, packageVersion)

	// load package manifest
	pkg, err := manager.loadPackage(packageVersionManifestPath)
	if err != nil {
		return nil, err
	}

	// validate?

	// interpolate package
	pkg, err = manager.interpolatePackageManifest(pkg, map[string]string{
		"/version" : packageVersion,
	})
	if err != nil {
		return nil, err
	}

	packageContext := NewContext(root, bin, packagesRoot, packagePath, packageVersionPath, packageVersionManifestPath)

	// turn package manifest into packages.Package
	// return package
	return manager.convertManifestToPackage(pkg, packageContext)
}

func (manager *manager) getPackagePath(packagesRoot, packageName string) (string, error) {
	if strings.TrimSpace(packageName) == "" {
		return "", fmt.Errorf("package name is required")
	}

	packagePath := filepath.Join(packagesRoot, packageName)
	return packagePath, nil
}

func (manager *manager) getPackageVersion(packagePath, packageVersion string) (string, error) {
	useLatestVersion := len(strings.TrimSpace(packageVersion)) == 0

	if !useLatestVersion {
		return packageVersion, nil
	}

	return manager.findLatestPackageVersion(packagePath)
}

func (manager *manager) findLatestPackageVersion(packagePath string) (string, error) {
	files, err := afero.ReadDir(manager.fileSystem, packagePath)
	if err != nil {
		return "", err
	}
		
	var latest *semver.Version

	for _, file := range files {

		if !file.IsDir() {
			continue
		}

		version := file.Name()
		v, err := semver.NewVersion(version)
		if err != nil{
			return "", err
		}

		if latest == nil{
			latest = v
			continue
		}

		if v.GreaterThan(latest){
			latest = v
		}
	}

	if latest == nil{
		return "", fmt.Errorf("unable to determine latest version of package")
	}

	return latest.String(), nil
}

func (manager *manager) loadPackage(packageManifestPath string) (interface{}, error) {	
	loader := config.NewLoader(manager.fileSystem)
	return loader.LoadPackageAsInterface(packageManifestPath)
}

func (manager *manager) interpolatePackageManifest(pkg interface{}, values map[string]string) (interface{}, error) {

	template := templates.NewTemplate(pkg)
	dictionary := collections.NewDictionaryFromMap(values)
	resolver := templates.NewDictionaryResolver(dictionary)

	return template.Evaluate(resolver)	
}

func (manager *manager) convertManifestToPackage(manifest interface{}, packageContext PackageContext) (Package, error) {
	pkg := &config.Package{}

	// convert to config structure
	err := mapstructure.Decode(manifest, pkg)
	if err != nil{
		return nil, err
	}

	// convert task list
	taskList := []tasks.Task{}
	for _, target := range pkg.Targets {
		for _, task := range target.Tasks{
			tsk,err := manager.taskProviders.Decode(task)
			if err !=nil{
				return nil, err
			}
			taskList = append(taskList, tsk)
		}
	}

	// convert package metadata
	return New(pkg.Name, pkg.Version, packageContext, taskList...), nil	
}
