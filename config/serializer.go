package config

import (
	yaml "gopkg.in/yaml.v2"
)

func Serialize(content []byte) (*Config, error) {
	config := &Config{}
	err := yaml.UnmarshalStrict([]byte(content), config)
	return config, err
}

func SerializeString(content string) (*Config, error) {
	return Serialize([]byte(content))
}
