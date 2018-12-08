package config

// Config represents a grouping of environments, stores and packages
type Config struct {
	Stores    []Store            `yaml:"stores"`
	Processes []Process          `yaml:"processes"`
	Imports   []PackageReference `yaml:"imports"`
}

// Store represents a configuration store
type Store struct {
	Name      string   `yaml:"name"`
	StoreType string   `yaml:"type"`
	Stores    []string `yaml:"stores"`

	Params map[string]string `yaml:"params"`
}

// Process represents a process under the given environment
type Process struct {
	Name   string            `yaml:"name"`
	Stores []string          `yaml:"stores"`
	Path   string            `yaml:"path"`
	Args   []string          `yaml:"args"`
	Vars   map[string]string `yaml:"env"`
}

// PackageReference represents a reference to a package
type PackageReference struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Package represents a versioned artifiact
type Package struct {
	Details PackageDetails `yaml:"package"`
}

// PackageDetails contains the details for the package
type PackageDetails struct {
	Name      string     `yaml:"name"`
	Version   string     `yaml:"version"`
	Platforms []Platform `yaml:"platforms"`
}

// Platform represents a package platform install instructions
type Platform struct {
	Name  string `yaml:"name"`
	Tasks []Task `yaml:"tasks"`
}

// Task represents a generic task to perform on the platform
type Task struct {
	Name   string                 `yaml:"name"`
	Type   string                 `yaml:"type"`
	Params map[string]interface{} `yaml:"params"`
}
