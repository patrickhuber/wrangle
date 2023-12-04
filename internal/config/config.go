package config

const (
	TagLatest = "latest"

	ConfigApiVersion = "wrangle/config/v1"
)

type Config struct {
	ApiVersion string            `json:"apiVersion" yaml:"apiVersion"`
	Metadata   map[string]string `json:"metadata" yaml:"metadata"`
	Spec       Spec              `json:"spec" yaml:"spec"`
}

type Spec struct {
	Feeds       []Feed            `json:"feeds" yaml:"feeds"`
	Stores      []Store           `json:"stores" yaml:"stores"`
	Environment map[string]string `json:"env" yaml:"env"`
	Packages    []Package         `json:"packages" yaml:"packages"`
}

type Feed struct {
	Name string `json:"name" yaml:"name"`
	Type string `json:"type" yaml:"type"`
	URI  string `json:"uri" yaml:"uri"`
}

type Store struct{}
type Environment struct{}
type Package struct {
	Name    string
	Version string
}
