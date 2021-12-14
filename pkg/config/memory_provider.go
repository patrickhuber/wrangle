package config

type memoryProvider struct {
	cfg *Config
}

func NewMemoryReader(cfg *Config) Reader {
	return &memoryProvider{
		cfg: cfg,
	}
}

func (p *memoryProvider) Get() (*Config, error) {
	return p.cfg, nil
}

func (p *memoryProvider) Set(cfg *Config) error {
	p.cfg = cfg
	return nil
}
