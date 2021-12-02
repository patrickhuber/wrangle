package memory

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type packageVersionRepository struct {
	items map[string]*feed.Item
}

func (r *packageVersionRepository) Get(packageName string, version string) (*packages.PackageVersion, error) {
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

func (r *packageVersionRepository) List(packageName string, latestVersion string, query *feed.ItemReadExpandPackage) ([]*packages.PackageVersion, error) {
	item, ok := r.items[packageName]
	versions := []*packages.PackageVersion{}

	if !ok {
		return versions, nil
	}

	for _, v := range item.Package.Versions {
		isMatch := query.IsMatch(v.Version, latestVersion)
		if isMatch {
			versions = append(versions, v)
		}
	}
	return versions, nil
}
