package packages

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
	Properties map[string]any
}
