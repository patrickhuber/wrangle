package packages

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"
)

type packageContext struct {
	root                       string
	bin                        string
	packagesRoot               string
	packagePath                string
	packageVersionPath         string
	packageVersionManifestPath string
}

// PackageContext defines the context for the given package
type PackageContext interface {
	Root() string
	Bin() string
	PackagesRoot() string
	PackagePath() string
	PackageVersionPath() string
	PackageVersionManifestPath() string
}

// NewContext creates a new package context with the given parameters
func NewContext(root, bin, packagesRoot, packagePath, packageVersionPath, packageVersionManifestPath string) PackageContext {
	return &packageContext{
		root:                       root,
		bin:                        bin,
		packagesRoot:               packagesRoot,
		packagePath:                packagePath,
		packageVersionPath:         packageVersionPath,
		packageVersionManifestPath: packageVersionManifestPath,
	}
}

// NewDefaultContext context sets the context to the default paths
func NewDefaultContext(root, packageName, packageVersion string) PackageContext {
	bin := filepath.Join(root, "/bin")
	packagesRoot := filepath.Join(root, "/packages")
	packagePath := filepath.Join(packagesRoot, packageName)
	packageVersionPath := filepath.Join(packagePath, packageVersion)
	packageVersionManifestPath := filepath.Join(packageVersionPath, fmt.Sprintf("%s.%s.yml", packageName, packageVersion))
	return NewContext(root, bin, packagesRoot, packagePath, packageVersionPath, packageVersionManifestPath)
}

func (pc *packageContext) Root() string {
	return pc.root
}

func (pc *packageContext) Bin() string {
	return pc.bin
}

func (pc *packageContext) PackagesRoot() string {
	return pc.packagesRoot
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
