package feed

import "github.com/patrickhuber/wrangle/internal/packages"

type Item struct {
	Package   *packages.Package
	State     *State
	Template  string
	Platforms []*Platform
}

type State struct {
	LatestVersion string `mapstructure:"version" yaml:"version" json:"version"`
}

type Platforms struct {
	Platforms []*Platform
}

type Platform struct {
	Name          string
	Architectures []string
}
