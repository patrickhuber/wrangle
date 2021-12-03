package tasks

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
)

// Metadata defines the context for a package
type Metadata struct {
	Name                       string
	Version                    string
	PackagePath                string
	PackageVersionPath         string
	PackageVersionManifestPath string
}

// NewDefaultMetadata creates a new default context given the package parameters
func NewDefaultMetadata(cfg *config.Config, name, version string) *Metadata {

	packagePath := crosspath.Join(cfg.PackagePath, name)
	packageVersionPath := crosspath.Join(packagePath, version)
	packageManifestName := "package.yml"
	packageVersionManifestPath := crosspath.Join(packageVersionPath, packageManifestName)
	return &Metadata{
		Name:                       name,
		Version:                    version,
		PackagePath:                packagePath,
		PackageVersionPath:         packageVersionPath,
		PackageVersionManifestPath: packageVersionManifestPath,
	}
}
