package config

type defaultProvider struct {
	defaultConfig *Config
	provider      Provider
}

func NewDefaultableProvider(provider Provider, defaultConfig *Config) Provider {
	return &defaultProvider{
		provider:      provider,
		defaultConfig: defaultConfig,
	}
}

func (p *defaultProvider) Get() (*Config, error) {
	cfg, err := p.provider.Get()
	if err != nil {
		return p.defaultConfig, nil
	}
	return cfg, err
}

func (p *defaultProvider) Lookup() (*Config, bool, error) {
	cfg, ok, err := p.provider.Lookup()
	if err != nil {
		return nil, false, err
	}
	if !ok {
		return p.defaultConfig, true, nil
	}
	return cfg, ok, err
}

func (p *defaultProvider) Set(cfg *Config) error {
	return p.provider.Set(cfg)
}
