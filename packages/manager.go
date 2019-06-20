package packages

import (
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/wrangle/collections"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/tasks"
	"github.com/patrickhuber/wrangle/templates"
)

type manager struct {
	fileSystem       filesystem.FsWrapper
	taskProviders    tasks.ProviderRegistry
	contextProvider  ContextProvider
	manifestProvider ManifestProvider
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
	contextProvider := NewFsContextProvider(manager.fileSystem, root, bin, packagesRoot)
	packageContext, err := contextProvider.Get(packageName, packageVersion)
	if err != nil {
		return nil, err
	}

	manifestProvider := NewFsManifestProvider(manager.fileSystem, packagesRoot)
	manifest, err := manifestProvider.GetInterface(packageContext)
	if err != nil {
		return nil, err
	}

	// validate?

	// interpolate package
	manifest, err = manager.interpolatePackageManifest(manifest, map[string]string{
		"/version": packageVersion,
	})
	if err != nil {
		return nil, err
	}

	// turn package manifest into packages.Package
	// return package
	return manager.convertManifestToPackage(manifest, packageContext)
}

func (manager *manager) interpolatePackageManifest(pkg interface{}, values map[string]string) (interface{}, error) {

	template := templates.NewTemplate(pkg)
	dictionary := collections.NewDictionaryFromMap(values)
	resolver := templates.NewDictionaryResolver(dictionary)

	return template.Evaluate(resolver)
}

func (manager *manager) convertManifestToPackage(manifest interface{}, packageContext PackageContext) (Package, error) {
	pkg := &Manifest{}

	// convert to config structure
	err := mapstructure.Decode(manifest, pkg)
	if err != nil {
		return nil, err
	}

	// convert task list
	taskList := []tasks.Task{}
	for _, target := range pkg.Targets {
		for _, task := range target.Tasks {
			tsk, err := manager.taskProviders.Decode(task)
			if err != nil {
				return nil, err
			}
			taskList = append(taskList, tsk)
		}
	}

	// convert package metadata
	return New(pkg.Name, pkg.Version, packageContext, taskList...), nil
}
