package config

import "fmt"

type memory struct {
	cfg *Config
}

func NewMemoryProvider() Provider {
	return &memory{}
}

func (m *memory) Get() (*Config, error) {
	if m.cfg == nil {
		return nil, fmt.Errorf("configuration is nil")
	}
	return m.cfg, nil
}

func (m *memory) Lookup() (*Config, bool, error) {
	if m.cfg == nil {
		return nil, false, nil
	}
	return m.cfg, true, nil
}

func (m *memory) Set(cfg *Config) error {
	m.cfg = cfg
	return nil
}
