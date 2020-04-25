package packages

import (
	"fmt"

	"github.com/patrickhuber/wrangle/tasks"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/settings"
)

type packageContext struct {
	paths                      *settings.Paths
	packagePath                string
	packageVersionPath         string
	packageVersionManifestPath string
	variables                  map[string]interface{}
}

// PackageContext defines the context for the given package
type PackageContext interface {
	tasks.TaskContext
}

// NewContext creates a new package context with the given parameters
func NewContext(paths *settings.Paths, packagePath, packageVersionPath, packageVersionManifestPath string) PackageContext {
	return &packageContext{
		paths:                      paths,
		packagePath:                packagePath,
		packageVersionPath:         packageVersionPath,
		packageVersionManifestPath: packageVersionManifestPath,
	}
}

// NewDefaultContext context sets the context to the default paths
func NewDefaultContext(root, packageName, packageVersion string) PackageContext {
	bin := filepath.Join(root, "bin")
	packagesRoot := filepath.Join(root, "packages")
	paths := &settings.Paths{
		Root:     root,
		Bin:      bin,
		Packages: packagesRoot,
	}
	packagePath := filepath.Join(packagesRoot, packageName)
	packageVersionPath := filepath.Join(packagePath, packageVersion)
	packageVersionManifestPath := filepath.Join(packageVersionPath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
	return NewContext(paths, packagePath, packageVersionPath, packageVersionManifestPath)
}

func (pc *packageContext) Root() string {
	return pc.paths.Root
}

func (pc *packageContext) Bin() string {
	return pc.paths.Bin
}

func (pc *packageContext) PackagesRoot() string {
	return pc.paths.Packages
}

func (pc *packageContext) PackagePath() string {
	return pc.packagePath
}

func (pc *packageContext) PackageVersionPath() string {
	return pc.packageVersionPath
}

func (pc *packageContext) PackageVersionManifestPath() string {
	return pc.packageVersionManifestPath
}

func (pc *packageContext) Variables() map[string]interface{} {
	return pc.variables
}

func (pc *packageContext) Variable(name string) (interface{}, bool) {
	value, ok := pc.variables[name]
	return value, ok
}

func (pc *packageContext) SetVariable(name string, value interface{}) {
	pc.variables[name] = value
}
