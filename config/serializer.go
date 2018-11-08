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

func SerializePackage(pkg *Package) (string, error) {
	data, err := yaml.Marshal(pkg)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
