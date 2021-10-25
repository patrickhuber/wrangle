package tasks

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

// Context defines the context for a package
type Context struct {
	Name                       string
	Version                    string
	PackagePath                string
	PackageVersionPath         string
	PackageVersionManifestPath string
}

// NewDefaultContext creates a new default context given the package parameters
func NewDefaultContext(cfg *config.Config, name, version string) *Context {

	packagePath := crosspath.Join(cfg.PackagePath, name)
	packageVersionPath := crosspath.Join(packagePath, version)
	packageManifestName := fmt.Sprintf("%s.%s.yml", name, version)
	packageVersionManifestPath := crosspath.Join(packageVersionPath, packageManifestName)
	return &Context{
		Name:                       name,
		Version:                    version,
		PackagePath:                packagePath,
		PackageVersionPath:         packageVersionPath,
		PackageVersionManifestPath: packageVersionManifestPath,
	}
}
