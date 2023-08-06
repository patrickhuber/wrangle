package actions

import (
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/wrangle/pkg/config"
)

// Metadata defines the context for a package
type Metadata struct {
	Name                       string
	Version                    string
	PackagePath                string
	PackageVersionPath         string
	PackageVersionManifestPath string
}

type metadataProvider struct {
	path *filepath.Processor
}

type MetadataProvider interface {
	Get(cfg *config.Config, name, version string) *Metadata
}

func NewMetadataProvider(path *filepath.Processor) MetadataProvider {
	return &metadataProvider{
		path: path,
	}
}

func (p *metadataProvider) Get(cfg *config.Config, name, version string) *Metadata {
	packagePath := p.path.Join(cfg.Paths.Packages, name)
	packageVersionPath := p.path.Join(packagePath, version)
	packageManifestName := "package.yml"
	packageVersionManifestPath := p.path.Join(packageVersionPath, packageManifestName)
	return &Metadata{
		Name:                       name,
		Version:                    version,
		PackagePath:                packagePath,
		PackageVersionPath:         packageVersionPath,
		PackageVersionManifestPath: packageVersionManifestPath,
	}
}
