package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

type PackageVersionRepository interface {
	Get(packageName string, version string) (*packages.PackageVersion, error)
	List(packageName string, latestVersion string, query *ItemReadExpandPackage) ([]*packages.PackageVersion, error)
}
