package config

import (
	yaml "gopkg.in/yaml.v2"
)

func DeserializeConfig(content []byte) (*Config, error) {
	config := &Config{}
	err := yaml.UnmarshalStrict(content, config)
	return config, err
}

func DeserializeConfigString(content string) (*Config, error) {
	return DeserializeConfig([]byte(content))
}

func DeserializePackage(content []byte) (*Package, error) {
	pkg := &Package{}
	err := yaml.UnmarshalStrict(content, pkg)
	return pkg, err
}

func DeserializePackageString(content string) (*Package, error) {
	return DeserializePackage([]byte(content))
}

// DeserializePackageAsInterface returns an interface representation of the package as a map of interfaces
func DeserializePackageAsInterface(content []byte) (interface{}, error) {
	m := make(map[interface{}]interface{})
	err := yaml.Unmarshal(content, m)
	return m, err
}

func SerializePackage(pkg *Package) (string, error) {
	data, err := yaml.Marshal(pkg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func SerializePackageFromInterface(data interface{}) (string, error) {
	dataBytes, err := yaml.Marshal(data)
	if err != nil {
		return "", err
	}
	return string(dataBytes), nil
}
