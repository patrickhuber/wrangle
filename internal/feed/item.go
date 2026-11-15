package feed

import "github.com/patrickhuber/wrangle/internal/packages"

type Item struct {
	Package   *packages.Package
	State     *State
	Template  string
	Platforms []*Platform
	Resource  *ResourceConfig
}

type State struct {
	LatestVersion string `mapstructure:"version" yaml:"version" json:"version"`
}

type Platforms struct {
	Platforms []*Platform
}

type Platform struct {
	Name          string   `yaml:"platform" mapstructure:"platform"`
	Architectures []string `yaml:"architectures" mapstructure:"architectures"`
}

// ResourceConfig holds the resource configuration for a feed item (resource.yml)
type ResourceConfig struct {
	Name   string            `yaml:"name"`
	Type   string            `yaml:"type"`
	Source map[string]string `yaml:"source"`
}

// ItemResourceFile wraps ResourceConfig matching the resource.yml file format
type ItemResourceFile struct {
	Resource *ResourceConfig `yaml:"resource"`
}
