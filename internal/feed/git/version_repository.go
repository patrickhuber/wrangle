package git

import (
	"fmt"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/util"
	"github.com/mitchellh/mapstructure"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
	"gopkg.in/yaml.v3"
)

type versionRepository struct {
	fs               billy.Filesystem
	workingDirectory string
	logger           log.Logger
	path             *filepath.Processor
}

func NewVersionRepository(fs billy.Filesystem, logger log.Logger, path *filepath.Processor, workingDirectory string) feed.VersionRepository {
	return &versionRepository{
		fs:               fs,
		workingDirectory: workingDirectory,
		logger:           logger,
		path:             path,
	}
}

func (s *versionRepository) List(name string) ([]*packages.Version, error) {
	s.logger.Tracef("versionRepository.List %s", name)
	packageDirectory := s.path.Join(s.workingDirectory, name)
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
	s.logger.Tracef("versionRepository.Get %s@%s", name, version)
	manifest, err := s.GetManifest(name, version)
	if err != nil {
		return nil, err
	}
	v := &packages.Version{
		Version:  manifest.Package.Version,
		Manifest: manifest,
	}
	return v, nil
}

func (s *versionRepository) GetManifest(name string, version string) (*packages.Manifest, error) {
	s.logger.Tracef("versionRepository.GetManifest %s@%s", name, version)
	manifestPath := s.path.Join(s.workingDirectory, name, version, "package.yml")
	content, err := util.ReadFile(s.fs, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("%w %s", err, manifestPath)
	}

	// validate with mapstructure package
	// convert to object
	var obj any
	err = yaml.Unmarshal(content, &obj)
	if err != nil {
		return nil, err
	}

	// then map to struct and validate
	manifest := &packages.Manifest{}
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:      manifest,
		ErrorUnused: true,
		ErrorUnset:  true,
	})
	if err != nil {
		return nil, err
	}

	err = decoder.Decode(obj)
	if err != nil {
		return nil, err
	}

	return manifest, err
}

func (s *versionRepository) Save(name string, version *packages.Version) error {
	s.logger.Tracef("versionRepository.Save %s@%s", name, version.Version)

	manifest := version.Manifest

	versionPath := s.path.Join(s.workingDirectory, name, version.Version)
	err := s.fs.MkdirAll(versionPath, 0600)
	if err != nil {
		return err
	}

	content, err := yaml.Marshal(manifest)
	if err != nil {
		return err
	}

	versionPackagePath := s.path.Join(versionPath, "package.yml")
	return util.WriteFile(s.fs, versionPackagePath, content, 0644)
}

func (s *versionRepository) Remove(name string, version string) error {
	path := s.path.Join(s.workingDirectory, name, version)
	return s.fs.Remove(path)
}
