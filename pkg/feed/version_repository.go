package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

type VersionRepository interface {
	Get(packageName string, version string) (*packages.PackageVersion, error)
	List(packageName string, query *ItemReadExpandPackage) ([]*packages.PackageVersion, error)
	Update(packageName string, command *VersionUpdate) ([]*packages.PackageVersion, error)
}
