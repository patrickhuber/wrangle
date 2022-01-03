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
	item, ok := r.items[name]
	if !ok {
		return versions, nil
	}
	versions = append(versions, item.Package.Versions...)
	return versions, nil
}

func (r *versionRepository) Save(name string, version *packages.Version) error {
	item, ok := r.items[name]
	if !ok {
		return nil
	}
	saved := false
	for i, v := range item.Package.Versions {
		if v.Version == version.Version {
			item.Package.Versions[i] = v
			saved = true
			break
		}
	}
	if !saved {
		item.Package.Versions = append(item.Package.Versions, version)
	}
	return nil
}

func (r *versionRepository) Remove(name string, version string) error {
	item, ok := r.items[name]
	if !ok {
		return nil
	}
	index := -1
	for i, v := range item.Package.Versions {
		if v.Version == version {
			index = i
			break
		}
	}
	if index != -1 {
		item.Package.Versions = append(item.Package.Versions[:index], item.Package.Versions[index+1:]...)
	}
	return nil
}
