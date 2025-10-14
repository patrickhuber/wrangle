package config

import (
	"github.com/patrickhuber/go-config"
	"github.com/patrickhuber/wrangle/internal/global"
)

type CliProvider struct {
	ctx     CliContext
	flagMap map[string]string
}

func NewCliProvider(ctx CliContext) config.Provider {
	flagMap := map[string]string{
		global.FlagBin:          global.EnvBin,
		global.FlagSystemConfig: global.EnvSystemConfig,
		global.FlagUserConfig:   global.EnvUserConfig,
		global.FlagRoot:         global.EnvRoot,
		global.FlagLogLevel:     global.EnvLogLevel,
	}

	return &CliProvider{
		ctx:     ctx,
		flagMap: flagMap,
	}
}

func (p *CliProvider) Get(ctx *config.GetContext) (any, error) {
	m := map[string]string{}
	for f, e := range p.flagMap {
		// only set the flag if the value is not empty
		value := p.ctx.String(f)
		if value != "" {
			m[e] = value
		}
	}
	return map[string]any{
		"spec": map[string]any{
			"env": m,
		},
	}, nil
}
