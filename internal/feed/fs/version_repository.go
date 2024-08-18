package fs

import (
	"fmt"

	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
	"gopkg.in/yaml.v3"
)

type versionRepository struct {
	fs               fs.FS
	workingDirectory string
	path             *filepath.Processor
}

func NewVersionRepository(fs fs.FS, path *filepath.Processor, workingDirectory string) feed.VersionRepository {
	return &versionRepository{
		fs:               fs,
		path:             path,
		workingDirectory: workingDirectory,
	}
}

func (r *versionRepository) Save(name string, version *packages.Version) error {
	versionPath := r.GetVersionFolderPath(name, version.Version)
	err := r.fs.MkdirAll(versionPath, 0644)
	if err != nil {
		return err
	}
	versionFile := r.GetVersionFilePath(name, version.Version)
	manifest := version.Manifest

	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	return r.fs.WriteFile(versionFile, data, 0644)
}

func (r *versionRepository) Get(name string, version string) (*packages.Version, error) {
	versionFile := r.GetVersionFilePath(name, version)
	data, err := r.fs.ReadFile(versionFile)
	if err != nil {
		return nil, err
	}
	manifest := &packages.Manifest{}
	err = yaml.Unmarshal(data, manifest)
	if err != nil {
		return nil, err
	}
	if manifest.Package == nil {
		return nil, fmt.Errorf("invalid package %s. manifest.Package is nil", versionFile)
	}
	v := &packages.Version{
		Version:  manifest.Package.Version,
		Manifest: manifest,
	}
	return v, nil
}

func (r *versionRepository) List(name string) ([]*packages.Version, error) {
	packagePath := r.GetPackageFolderPath(name)
	files, err := r.fs.ReadDir(packagePath)
	if err != nil {
		return nil, err
	}
	versions := []*packages.Version{}
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		version, err := r.Get(name, file.Name())
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (r *versionRepository) Remove(name string, version string) error {
	versionPath := r.path.Join(r.workingDirectory, name, version)
	return r.fs.RemoveAll(versionPath)
}

func (r *versionRepository) GetPackageFolderPath(name string) string {
	return r.path.Join(r.workingDirectory, name)
}

func (r *versionRepository) GetVersionFilePath(name, version string) string {
	return r.path.Join(r.GetVersionFolderPath(name, version), "package.yml")
}

func (r *versionRepository) GetVersionFolderPath(name, version string) string {
	return r.path.Join(r.GetPackageFolderPath(name), version)
}
