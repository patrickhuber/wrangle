package config

import (
	"os"
	"reflect"

	"github.com/caarlos0/env/v6"
)

type envProvider struct {
}

// NewEnvProvider creates a new environment variable provider for config
func NewEnvProvider() Provider {
	return &envProvider{}
}

func (p *envProvider) Get() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (p *envProvider) Set(config *Config) error {
	t := reflect.TypeOf(config)
	v := reflect.ValueOf(config)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("env")
		if tag == "" {
			continue
		}
		value := v.Elem().String()
		err := os.Setenv(field.Name, value)
		if err != nil {
			return err
		}
	}
	return nil
}
