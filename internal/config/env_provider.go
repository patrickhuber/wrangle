package config

import (
	"strings"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/wrangle/internal/global"
)

type EnvProvider struct {
	environment env.Environment
	prefixes    []string
}

func NewEnvProvider(environment env.Environment, prefixes ...string) config.Provider {
	if len(prefixes) == 0 {
		prefixes = append(prefixes, global.EnvPrefix)
	}
	return &EnvProvider{
		environment: environment,
		prefixes:    prefixes,
	}
}

func (p *EnvProvider) Get(ctx *config.GetContext) (any, error) {
	// use the environment to get the configuration
	m := p.environment.Export()
	envConfig := map[string]any{}
	for k, v := range m {
		// check if the key starts with any of the prefixes
		for _, prefix := range p.prefixes {
			if strings.HasPrefix(k, prefix) {
				envConfig[k] = v
			}
		}
	}

	// return the env config as is
	return map[string]any{
		"spec": map[string]any{
			"env": envConfig,
		},
	}, nil
}
