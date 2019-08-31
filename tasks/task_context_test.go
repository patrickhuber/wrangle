package tasks_test

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/tasks"
)

type fakeTaskContext struct {
	root                       string
	bin                        string
	packagesRoot               string
	packagePath                string
	packageVersionPath         string
	packageVersionManifestPath string
	variables                  map[string]interface{}
}

func newFakeTaskContext(root, packageName, packageVersion string) tasks.TaskContext {

	bin := filepath.Join(root, "/bin")
	packagesRoot := filepath.Join(root, "/packages")
	packagePath := filepath.Join(packagesRoot, packageName)
	packageVersionPath := filepath.Join(packagePath, packageVersion)
	packageVersionManifestPath := filepath.Join(packageVersionPath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
	variables := make(map[string]interface{})
	return &fakeTaskContext{
		bin:                        bin,
		root:                       root,
		packagesRoot:               packagesRoot,
		packagePath:                packagePath,
		packageVersionPath:         packageVersionPath,
		packageVersionManifestPath: packageVersionManifestPath,
		variables:                  variables,
	}
}

func (tc *fakeTaskContext) Bin() string {
	return tc.bin
}

func (tc *fakeTaskContext) Root() string {
	return tc.root
}

func (tc *fakeTaskContext) PackagesRoot() string {
	return tc.packagesRoot
}

func (tc *fakeTaskContext) PackagePath() string {
	return tc.packagePath
}

func (tc *fakeTaskContext) PackageVersionPath() string {
	return tc.packageVersionPath
}

func (tc *fakeTaskContext) PackageVersionManifestPath() string {
	return tc.packageVersionManifestPath
}

func (tc *fakeTaskContext) Variables() map[string]interface{} {
	return tc.variables
}

func (tc *fakeTaskContext) Variable(name string) (interface{}, bool) {
	value, ok := tc.variables[name]
	return value, ok
}

func (tc *fakeTaskContext) SetVariable(name string, value interface{}) {
	tc.variables[name] = value
}
