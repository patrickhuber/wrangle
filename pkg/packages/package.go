package packages

import (
	"gopkg.in/yaml.v2"
)

// Package defines a software package and how it is installed
type Package struct {
	Name     string
	Versions []*PackageVersion
}

// PackageVersion returns the version of a package
type PackageVersion struct {
	Version string
	Targets []*PackageTarget
}

// PackageTarget defines a target architecture and platform for the series of tasks to run
type PackageTarget struct {
	Platform     string
	Architecture string
	Tasks        []*PackageTargetTask
}

// PackageTargetTask defines a target task to run for the given target
type PackageTargetTask struct {
	Name       string
	Properties map[string]string
}

func FromYaml(manifest string) (*Package, error) {
	var manifestStruct Manifest
	err := yaml.Unmarshal([]byte(manifest), &manifestStruct)
	if err != nil {
		return nil, err
	}
	return toPackage(&manifestStruct)
}

func toPackage(manifest *Manifest) (*Package, error) {
	return &Package{
		Name: manifest.Package.Name,
		Versions: []*PackageVersion{
			{
				Version: manifest.Package.Version,
				Targets: manifest.Package.Targets,
			},
		},
	}, nil
}
