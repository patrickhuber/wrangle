package git

import (
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
)

type versionRepository struct {
	fs               billy.Filesystem
	workingDirectory string
}

func NewVersionRepository(fs billy.Filesystem, workingDirectory string) feed.VersionRepository {
	return &versionRepository{
		fs:               fs,
		workingDirectory: workingDirectory,
	}
}

func (s *versionRepository) List(name string) ([]*packages.Version, error) {
	packageDirectory := crosspath.Join(s.workingDirectory, name)
	files, err := s.fs.ReadDir(packageDirectory)
	if err != nil {
		return nil, err
	}
	versions := []*packages.Version{}
	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		version, err := s.Get(name, f.Name())
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (s *versionRepository) Get(name string, version string) (*packages.Version, error) {
	manifest, err := s.GetManifest(name, version)
	if err != nil {
		return nil, err
	}
	return &packages.Version{
		Version: manifest.Package.Version,
		Targets: manifest.Package.Targets,
	}, nil
}

func (s *versionRepository) GetManifest(name string, version string) (*packages.Manifest, error) {
	manifestPath := crosspath.Join(s.workingDirectory, name, version, "package.yml")
	content, err := util.ReadFile(s.fs, manifestPath)
	if err != nil {
		return nil, err
	}
	manifest := &packages.Manifest{}
	err = yaml.Unmarshal(content, &manifest)
	return manifest, err
}

func (s *versionRepository) Save(name string, version *packages.Version) error {
	manifest := &packages.Manifest{
		Package: &packages.ManifestPackage{
			Name:    name,
			Version: version.Version,
			Targets: version.Targets,
		},
	}

	versionPath := crosspath.Join(s.workingDirectory, name, version.Version)
	err := s.fs.MkdirAll(versionPath, 0600)
	if err != nil {
		return err
	}

	content, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	versionPackagePath := crosspath.Join(versionPath, "package.yml")
	return util.WriteFile(s.fs, versionPackagePath, content, 0644)
}

func (s *versionRepository) Remove(name string, version string) error {
	path := crosspath.Join(s.workingDirectory, name, version)
	return s.fs.Remove(path)
}
