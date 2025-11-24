package config

const (
	TagLatest = "latest"

	ApiVersion = "wrangle/v1"
	Kind       = "Config"
)

type Config struct {
	ApiVersion string `json:"apiVersion" yaml:"apiVersion" toml:"apiVersion" mapstructure:"apiVersion"`
	Kind       string `json:"kind" yaml:"kind" toml:"kind" mapstructure:"kind"`
	Spec       Spec   `json:"spec" yaml:"spec" toml:"spec" mapstructure:"spec"`
}

type Spec struct {
	Feeds       []Feed            `json:"feeds" yaml:"feeds" toml:"feeds" mapstructure:"feeds"`
	Stores      []Store           `json:"stores" yaml:"stores" toml:"stores" mapstructure:"stores"`
	Environment map[string]string `json:"env" yaml:"env" toml:"env" mapstructure:"env"`
	Packages    []Package         `json:"packages" yaml:"packages" toml:"packages" mapstructure:"packages"`
}

type Feed struct {
	Name string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Type string `json:"type" yaml:"type" toml:"type" mapstructure:"type"`
	URI  string `json:"uri" yaml:"uri" toml:"uri" mapstructure:"uri"`
}

/*
A store for key vaule pairs

	type: azure.keyvault
	properties:
		uri: {key vault uri} // (required)

	type: github.secrets
	properties:
		org: {org name}	// (required)
		repo: {repo name} // (required)
*/
type Store struct {
	// the name of the store
	Name string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`

	// discriminator type (azure.keyvault, ...)
	Type string `json:"type" yaml:"type" toml:"type" mapstructure:"type"`

	// properties of the specific type
	Properties map[string]any `json:"properties" yaml:"properties" toml:"properties" mapstructure:"properties"`
}

type Package struct {
	Name    string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Version string `json:"version" yaml:"version" toml:"version" mapstructure:"version"`
}
