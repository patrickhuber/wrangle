package config

const (
	TagLatest = "latest"
)

// Config represents a configuration document for wrangle
type Config struct {
	Paths      *Paths       `yaml:"paths,omitempty"`
	Feeds      []*Feed      `yaml:"feeds,omitempty"`
	Stores     []*Store     `yaml:"stores,omitempty"`
	Processes  []*Process   `yaml:"processes,omitempty"`
	References []*Reference `yaml:"references,omitempty"`
}

// Paths contains the defualt paths for the config
type Paths struct {
	Packages string `yaml:"packages"`
	Bin      string `yaml:"bin"`
	Root     string `yaml:"root"`
}

// Feed represents a pacakge feed
type Feed struct {
	Name string `yaml:"name"`
	URI  string `yaml:"uri"`
	Type string `yaml:"type"`
}

// Store represents a configuration store
type Store struct {
	Name      string   `yaml:"name"`
	StoreType string   `yaml:"type"`
	Stores    []string `yaml:"stores"`

	Params map[string]string `yaml:"params"`
}

// Reference represents a reference to a package
type Reference struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

// Process represents a process under the given environment
type Process struct {
	Name   string            `yaml:"name"`
	Stores []string          `yaml:"stores"`
	Path   string            `yaml:"path"`
	Args   []string          `yaml:"args"`
	Vars   map[string]string `yaml:"env"`
}
