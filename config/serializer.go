package config

import (
	yaml "gopkg.in/yaml.v2"
)

func SerializeConfig(content []byte) (*Config, error) {
	config := &Config{}
	err := yaml.UnmarshalStrict(content, config)
	return config, err
}

func SerializeConfigString(content string) (*Config, error) {
	return SerializeConfig([]byte(content))
}

func SerializePackage(content []byte) (*Package, error) {
	pkg := &Package{}
	err := yaml.UnmarshalStrict(content, pkg)
	return pkg, err
}

func SerializePackageString(content string) (*Package, error) {
	return SerializePackage([]byte(content))
}
