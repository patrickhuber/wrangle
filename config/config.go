package config

// Config represents a grouping of environments and config sources
type Config struct {
	ConfigSources []ConfigSource `yaml:"config-sources"`
	Environments  []Environment  `yaml:"environments"`
}

// ConfigSource represents a source of configuration
type ConfigSource struct {
	Name             string            `yaml:"name"`
	ConfigSourceType string            `yaml:"type"`
	Config           string            `yaml:"config"`
	Params           map[string]string `yaml:"params"`
}

// Environment represents a grouping of processes
type Environment struct {
	Name      string    `yaml:"name"`
	Processes []Process `yaml:"processes"`
}

// Process represents a process under the given environment
type Process struct {
	Name   string            `yaml:"name"`
	Config string            `yaml:"config"`
	Path   string            `yaml:"path"`
	Args   []string          `yaml:"args"`
	Vars   map[string]string `yaml:"env"`
}
