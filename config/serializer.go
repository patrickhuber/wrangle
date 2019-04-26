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
