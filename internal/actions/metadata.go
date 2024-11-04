package actions

import (
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/global"
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
	path filepath.Provider
}

type MetadataProvider interface {
	Get(cfg *config.Config, name, version string) *Metadata
}

func NewMetadataProvider(path filepath.Provider) MetadataProvider {
	return &metadataProvider{
		path: path,
	}
}

func (p *metadataProvider) Get(cfg *config.Config, name, version string) *Metadata {
	envPackages := cfg.Spec.Environment[global.EnvPackages]
	packagePath := p.path.Join(envPackages, name)
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
