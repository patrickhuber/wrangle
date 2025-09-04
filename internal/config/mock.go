package config

type mock struct {
	cfg Config
}

func NewMock(cfg Config) Configuration {
	return &mock{
		cfg: cfg,
	}
}

func (m *mock) Get() (Config, error) {
	return m.cfg, nil
}
