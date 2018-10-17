package config

// Config represents a grouping of environments, stores and packages
type Config struct {
	Stores    []Store   `yaml:"stores"`
	Packages  []Package `yaml:"packages"`
	Processes []Process `yaml:"processes"`
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

// Package represents a versioned artifiact
type Package struct {
	Name      string     `yaml:"name"`
	Version   string     `yaml:"version"`
	Platforms []Platform `yaml:"platforms"`
}

// Platform represents a package platform install instructions
type Platform struct {
	Name     string    `yaml:"name"`
	Download *Download `yaml:"download"`
	Alias    string    `yaml:"alias"`
	Extract  *Extract  `yaml:"extract"`
}

// Download represents a package download
type Download struct {
	Out string `yaml:"out"`
	URL string `yaml:"url"`
}

// Extract represents package extract information
type Extract struct {
	Out    string `yaml:"out"`
	Filter string `yaml:"filter"`
}
