package config

// Config represents a grouping of environments and config sources
type Config struct {
	ConfigSources []ConfigSource `yaml:"config-sources"`
	Environments  []Environment  `yaml:"environments"`
	Packages      []Package      `yaml:"packages"`
}

// ConfigSource represents a source of configuration
type ConfigSource struct {
	Name             string            `yaml:"name"`
	ConfigSourceType string            `yaml:"type"`
	Configurations   []string          `yaml:"configurations"`
	Params           map[string]string `yaml:"params"`
}

// Environment represents a grouping of processes
type Environment struct {
	Name      string    `yaml:"name"`
	Processes []Process `yaml:"processes"`
}

// Process represents a process under the given environment
type Process struct {
	Name           string            `yaml:"name"`
	Configurations []string          `yaml:"configurations"`
	Path           string            `yaml:"path"`
	Args           []string          `yaml:"args"`
	Vars           map[string]string `yaml:"env"`
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
