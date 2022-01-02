package packages

import (
	"gopkg.in/yaml.v2"
)

// Package defines a software package and how it is installed
type Package struct {
	Name     string
	Versions []*Version
}

// Version returns the version of a package
type Version struct {
	Version string
	Targets []*Target
}

// Target defines a target architecture and platform for the series of tasks to run
type Target struct {
	Platform     string
	Architecture string
	Tasks        []*Task
}

// Task defines a target task to run for the given target
type Task struct {
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
		Versions: []*Version{
			{
				Version: manifest.Package.Version,
				Targets: manifest.Package.Targets,
			},
		},
	}, nil
}
