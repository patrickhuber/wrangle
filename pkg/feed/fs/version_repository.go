package fs

import (
	"os"

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

func (r *versionRepository) Get(packageName string, version string) (*packages.Version, error) {
	ok, err := r.fs.Exists(crosspath.Join(r.workingDirectory, packageName))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	ok, err = r.fs.Exists(crosspath.Join(r.workingDirectory, packageName, version))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	return nil, nil
}

func (r *versionRepository) List(packageName string, query *feed.ItemReadExpandPackage) ([]*packages.Version, error) {
	ok, err := r.fs.Exists(crosspath.Join(r.workingDirectory, packageName))
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, os.ErrNotExist
	}

	versions := []*packages.Version{}
	files, err := r.fs.ReadDir(crosspath.Join(r.workingDirectory, packageName))
	if err != nil {
		return nil, err
	}

	latestVersion, err := r.GetLatestVersion(packageName)
	if err != nil {
		return nil, err
	}

	for _, f := range files {
		if !f.IsDir() {
			continue
		}
		version := &packages.Version{
			Version: f.Name(),
		}
		if query == nil || query.IsMatch(f.Name(), latestVersion) {
			versions = append(versions, version)
		}
	}
	return versions, nil
}

func (r *versionRepository) GetLatestVersion(packageName string) (string, error) {
	latestVersion := ""
	bytes, err := r.fs.Read(crosspath.Join(r.workingDirectory, packageName, "state.yml"))
	if err != nil {
		return "", err
	}

	state := &feed.State{}
	err = yaml.Unmarshal(bytes, state)
	if err != nil {
		return "", err
	}
	if state != nil {
		latestVersion = state.LatestVersion
	}
	return latestVersion, nil
}

func (r *versionRepository) Update(packageName string, command *feed.VersionUpdate) ([]*packages.Version, error) {
	return nil, nil
}
