package packages

import "gopkg.in/yaml.v2"

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
	var data map[string]interface{}
	err := yaml.Unmarshal([]byte(manifest), data)
	if err != nil {
		return nil, err
	}
	return toPackage(data)
}

func toPackage(data map[string]interface{}) (*Package, error) {
	packageTargets := []*PackageTarget{}
	targets := data["targets"].([]map[string]interface{})
	for _, t := range targets {
		packageTargetTasks := []*PackageTargetTask{}
		tasks := t["tasks"].([]map[string]interface{})
		for _, tsk := range tasks {
			for name, value := range tsk {
				packageTargetTask := &PackageTargetTask{
					Name:       name,
					Properties: value.(map[string]string),
				}
				packageTargetTasks = append(packageTargetTasks, packageTargetTask)
			}
		}
		packageTarget := &PackageTarget{
			Platform:     t["platform"].(string),
			Architecture: t["architecture"].(string),
			Tasks:        packageTargetTasks,
		}
		packageTargets = append(packageTargets, packageTarget)
	}
	p := &Package{
		Name: data["name"].(string),
		Versions: []*PackageVersion{
			{
				Version: data["version"].(string),
				Targets: packageTargets,
			},
		},
	}
	return p, nil
}
