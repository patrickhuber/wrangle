package packages

func ManifestToPackageVersion(manifest *Manifest) *Version {
	targets := []*Target{}
	for _, tar := range manifest.Package.Targets {
		steps := []*Task{}
		for _, step := range tar.Steps {
			steps = append(steps, &Task{Name: step.Action, Properties: step.With})
		}
		targets = append(targets, &Target{
			Platform:     tar.Platform,
			Architecture: tar.Architecture,
			Tasks:        steps,
		})
	}
	return &Version{
		Version: manifest.Package.Version,
		Targets: targets,
	}
}

func PackageVersionToManifest(name string, version *Version) *Manifest {
	targets := []*ManifestTarget{}
	for _, tar := range version.Targets {
		steps := []ManifestStep{}
		for _, tsk := range tar.Tasks {
			step := ManifestStep{
				Action: tsk.Name,
				With:   tsk.Properties,
			}
			steps = append(steps, step)
		}
		targets = append(targets, &ManifestTarget{
			Platform:     tar.Platform,
			Architecture: tar.Architecture,
			Steps:        steps,
		})
	}
	return &Manifest{
		Package: &ManifestPackage{
			Name:    name,
			Version: version.Version,
			Targets: targets,
		},
	}
}
