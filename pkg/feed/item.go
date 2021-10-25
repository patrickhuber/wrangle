package feed

import "github.com/patrickhuber/wrangle/pkg/packages"

type Item struct {
	Package   *packages.Package
	State     *State
	Template  string
	Platforms []*Platform
}

type State struct {
	LatestVersion string
}

type Platform struct {
	Name          string
	Architectures []string
}
