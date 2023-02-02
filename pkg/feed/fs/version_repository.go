package fs

import (
	"fmt"

	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v3"
)

type versionRepository struct {
	fs               filesystem.FileSystem
	workingDirectory string
}

func NewVersionRepository(fs filesystem.FileSystem, workingDirectory string) feed.VersionRepository {
	return &versionRepository{
		fs:               fs,
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
	return r.fs.Write(versionFile, data, 0644)
}

func (r *versionRepository) Get(name string, version string) (*packages.Version, error) {
	versionFile := r.GetVersionFilePath(name, version)
	data, err := r.fs.Read(versionFile)
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
	versionPath := crosspath.Join(r.workingDirectory, name, version)
	return r.fs.RemoveAll(versionPath)
}

func (r *versionRepository) GetPackageFolderPath(name string) string {
	return crosspath.Join(r.workingDirectory, name)
}

func (r *versionRepository) GetVersionFilePath(name, version string) string {
	return crosspath.Join(r.GetVersionFolderPath(name, version), "package.yml")
}

func (r *versionRepository) GetVersionFolderPath(name, version string) string {
	return crosspath.Join(r.GetPackageFolderPath(name), version)
}
