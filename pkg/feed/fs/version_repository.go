package fs

import (
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
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
	versionPath := crosspath.Join(r.workingDirectory, name, version.Version)
	err := r.fs.MkdirAll(versionPath, 0644)
	if err != nil {
		return err
	}
	versionFile := crosspath.Join(versionPath, "package.yml")
	manifest := &packages.Manifest{
		Package: &packages.ManifestPackage{
			Name:    name,
			Version: version.Version,
			Targets: version.Targets,
		},
	}
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}
	return r.fs.Write(versionFile, data, 0644)
}

func (r *versionRepository) Get(name string, version string) (*packages.Version, error) {
	versionFile := crosspath.Join(r.workingDirectory, name, version, "package.yml")
	data, err := r.fs.Read(versionFile)
	if err != nil {
		return nil, err
	}
	v := &packages.Version{}
	err = yaml.Unmarshal(data, v)

	return v, err
}

func (r *versionRepository) List(name string) ([]*packages.Version, error) {
	packagePath := crosspath.Join(r.workingDirectory, name)
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
