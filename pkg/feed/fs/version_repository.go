package fs

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
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
	return nil
}
func (r *versionRepository) Get(name string, version string) (*packages.Version, error) {
	return nil, nil
}
func (r *versionRepository) List(name string) ([]*packages.Version, error) {
	return nil, nil
}
func (r *versionRepository) Remove(name string, version string) error {
	return nil
}
