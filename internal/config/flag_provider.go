package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type FlagProvider struct {
	provider config.Provider
}

func NewFlagProvider(args []string) config.Factory {
	flagToEnvMap := map[string]string{
		global.FlagBin:          global.EnvBin,
		global.FlagSystemConfig: global.EnvSystemConfig,
		global.FlagUserConfig:   global.EnvUserConfig,
		global.FlagRoot:         global.EnvRoot,
		global.FlagPackages:     global.EnvPackages,
		global.FlagLogLevel:     global.EnvLogLevel,
	}
	configFlags := []config.Flag{}
	for flag := range flagToEnvMap {
		configFlags = append(configFlags, &config.StringFlag{
			Name: flag,
		})
	}
	return config.NewFlag(
		configFlags,
		args,
		config.FlagOption{Transformers: []config.Transformer{
			config.FuncTypedTransformer(func(m map[string]any) (map[string]any, error) {
				result := map[string]any{}
				for flag, envVar := range flagToEnvMap {
					if value, ok := m[flag]; ok {
						result[envVar] = value
					}
				}
				return map[string]any{
					"spec": map[string]any{
						"env": result,
					},
				}, nil
			})}})
}
