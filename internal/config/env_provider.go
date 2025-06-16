package config

import (
	"strings"

	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/wrangle/internal/global"
)

type EnvProvider struct {
	environment env.Environment
}

func NewEnvProvider(environment env.Environment) config.Provider {
	return &EnvProvider{
		environment: environment,
	}
}

func (p *EnvProvider) Get(ctx *config.GetContext) (any, error) {
	// use the environment to get the configuration
	m := p.environment.Export()
	envConfig := map[string]any{}
	for k, v := range m {
		if strings.HasPrefix(k, global.EnvPrefix) {
			envConfig[k] = v
		}
	}

	// return the env config as is
	return map[string]any{
		"spec": map[string]any{
			"env": envConfig,
		},
	}, nil
}
