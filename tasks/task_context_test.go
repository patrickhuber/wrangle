package tasks_test

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/tasks"
)

type taskContext struct {
	root                       string
	bin                        string
	packagesRoot               string
	packagePath                string
	packageVersionPath         string
	packageVersionManifestPath string
}

func newTaskContext(root, packageName, packageVersion string) tasks.TaskContext {

	bin := filepath.Join(root, "/bin")
	packagesRoot := filepath.Join(root, "/packages")
	packagePath := filepath.Join(packagesRoot, packageName)
	packageVersionPath := filepath.Join(packagePath, packageVersion)
	packageVersionManifestPath := filepath.Join(packageVersionPath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))

	return &taskContext{
		bin:                        bin,
		root:                       root,
		packagesRoot:               packagesRoot,
		packagePath:                packagePath,
		packageVersionPath:         packageVersionPath,
		packageVersionManifestPath: packageVersionManifestPath,
	}
}

func (tc *taskContext) Bin() string {
	return tc.bin
}

func (tc *taskContext) Root() string {
	return tc.root
}

func (tc *taskContext) PackagesRoot() string {
	return tc.packagesRoot
}

func (tc *taskContext) PackagePath() string {
	return tc.packagePath
}

func (tc *taskContext) PackageVersionPath() string {
	return tc.packageVersionPath
}

func (tc *taskContext) PackageVersionManifestPath() string {
	return tc.packageVersionManifestPath
}
