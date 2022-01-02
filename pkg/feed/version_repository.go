package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

type VersionRepository interface {
	Save(name string, version *packages.Version) error
	Get(name string, version string) (*packages.Version, error)
	List(name string) ([]*packages.Version, error)
	Remove(name string, version string) error
}
