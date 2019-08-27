package packages

import "github.com/patrickhuber/wrangle/filesystem"

type fsManifestProvider struct {
	fs                filesystem.FileSystem
	packagesDirectory string
}

// NewFsManifestProvider creates a ManifestProvider for the given file system and package directory
func NewFsManifestProvider(fs filesystem.FileSystem, packagesDirectory string) ManifestProvider {
	return &fsManifestProvider{
		fs:                fs,
		packagesDirectory: packagesDirectory,
	}
}

func (p *fsManifestProvider) Get(context PackageContext) (*Manifest, error) {

	manifest, err := p.loadManifest(context.PackageVersionManifestPath())
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (p *fsManifestProvider) GetInterface(context PackageContext) (interface{}, error) {

	manifest, err := p.loadInterface(context.PackageVersionManifestPath())
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (p *fsManifestProvider) loadInterface(manifestPath string) (interface{}, error) {
	file, err := p.fs.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := NewYamlInterfaceReader(file)
	manifest, err := r.Read()
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func (p *fsManifestProvider) loadManifest(manifestPath string) (*Manifest, error) {
	file, err := p.fs.Open(manifestPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	r := NewYamlManifestReader(file)
	manifest, err := r.Read()
	if err != nil {
		return nil, err
	}

	return manifest, nil
}
