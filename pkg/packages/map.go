package packages

func ManifestToPackageVersion(manifest *Manifest) *Version {
	targets := []*Target{}
	for _, tar := range manifest.Package.Targets {
		tasks := []*Task{}
		for _, tsk := range tar.Tasks {
			for k, v := range tsk {
				tasks = append(tasks, &Task{
					Name:       k,
					Properties: v,
				})
			}
		}
		targets = append(targets, &Target{
			Platform:     tar.Platform,
			Architecture: tar.Architecture,
			Tasks:        tasks,
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
		tasks := []map[string]map[string]string{}
		for _, tsk := range tar.Tasks {
			task := map[string]map[string]string{}
			task[tsk.Name] = tsk.Properties
			tasks = append(tasks, task)
		}
		targets = append(targets, &ManifestTarget{
			Platform:     tar.Platform,
			Architecture: tar.Architecture,
			Tasks:        tasks,
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
