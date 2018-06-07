package config

type Config struct {
	ConfigSources []ConfigSource `yaml:"config-sources"`
	Processes     []Process      `yaml:"processes"`
}

type ConfigSource struct {
	Name             string            `yaml:"name"`
	ConfigSourceType string            `yaml:"type"`
	Config           string            `yaml:"config"`
	Params           map[string]string `yaml:"params"`
}

type Process struct {
	Name         string        `yaml:"name"`
	Environments []Environment `yaml:"environments"`
}

type Environment struct {
	Name    string            `yaml:"name"`
	Config  string            `yaml:"config"`
	Process string            `yaml:"process"`
	Args    []string          `yaml:"args"`
	Vars    map[string]string `yaml:"env"`
}
