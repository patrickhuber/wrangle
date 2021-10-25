package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

// Store defines methods for persistence of an item
type Store interface {
	SaveItem(item *Item)
	SavePackageVersion(packageName string, ver *packages.PackageVersion)
}
