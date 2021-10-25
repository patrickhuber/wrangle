package feed

import "github.com/patrickhuber/wrangle/pkg/config"

type manager struct {
	cfg *config.Config
}

// Manager provides feed aggregation management
type Manager interface {
	List() []*config.Feed
}

func NewManager(cfg *config.Config) Manager {
	return &manager{cfg: cfg}
}

func (m *manager) List() []*config.Feed {
	return m.cfg.Feeds
}
