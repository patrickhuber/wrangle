package memory

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type versionRepository struct {
	items map[string]*feed.Item
}

func NewVersionRepository(items map[string]*feed.Item) feed.VersionRepository {
	return &versionRepository{
		items: items,
	}
}

func (r *versionRepository) Get(packageName string, version string) (*packages.Version, error) {
	item, ok := r.items[packageName]
	if !ok {
		return nil, nil
	}
	for _, v := range item.Package.Versions {
		if v.Version == version {
			return v, nil
		}
	}
	return nil, nil
}

func (r *versionRepository) List(name string) ([]*packages.Version, error) {
	versions := []*packages.Version{}

	return versions, nil
}

func (r *versionRepository) Save(name string, version *packages.Version) error {
	return nil
}
func (r *versionRepository) Remove(name string, version string) error {
	return nil
}
