package config

// Config represents a configuration document for wrangle
type Config struct {
	PackagePath string
	BinPath     string
	RootPath    string

	Feeds     []*Feed             `yaml:"feeds"`
	Stores    []*Store            `yaml:"stores"`
	Processes []*Process          `yaml:"processes"`
	Imports   []*PackageReference `yaml:"imports"`
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

// PackageReference represents a reference to a package
type PackageReference struct {
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
