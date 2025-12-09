package config

const (
	TagLatest = "latest"

	ApiVersion = "wrangle/v1"
	Kind       = "Config"

	// Variable types based on BOSH variable types
	// https://bosh.io/docs/variable-types/
	VariableTypeCertificate = "certificate"
	VariableTypePassword    = "password"
	VariableTypeRSA         = "rsa"
	VariableTypeSSH         = "ssh"
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
	Variables   []Variable        `json:"variables" yaml:"variables" toml:"variables" mapstructure:"variables"`
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

	// dependencies are other stores this store requires to be evaluated first
	Dependencies []string `json:"dependencies" yaml:"dependencies" toml:"dependencies" mapstructure:"dependencies"`
}

type Package struct {
	Name    string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`
	Version string `json:"version" yaml:"version" toml:"version" mapstructure:"version"`
}

/*
A variable definition for generating secrets.
Variable values are not specified, only the properties that define how to generate them.

Supported types (based on BOSH variable types https://bosh.io/docs/variable-types/):

	type: certificate
	options:
		ca: {ca name}                      // name of a CA to sign the certificate (optional)
		common_name: {common name}         // common name for the certificate
		organization: {org}                // organization name (optional)
		alternative_names: [{alt names}]   // subject alternative names (optional)
		is_ca: {true|false}                // whether to generate a CA certificate (optional)
		extended_key_usage: [{usages}]     // extended key usage extensions (optional)
		duration: {days}                   // duration in days (optional)
		key_length: {length}               // key length, default 2048 (optional)

	type: password
	options:
		length: {length}                   // password length, default 20 (optional)

	type: rsa
	options:
		key_length: {length}               // key length, default 2048 (optional)

	type: ssh
	options:
		comment: {comment}                 // comment for the SSH key (optional)
*/
type Variable struct {
	// the name of the variable
	Name string `json:"name" yaml:"name" toml:"name" mapstructure:"name"`

	// discriminator type (certificate, password, rsa, ssh)
	Type string `json:"type" yaml:"type" toml:"type" mapstructure:"type"`

	// options for the specific variable type
	Options map[string]any `json:"options" yaml:"options" toml:"options" mapstructure:"options"`
}
