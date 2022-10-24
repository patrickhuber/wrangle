package memory

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

const ProviderType = "memory"

type provider struct {
	items  []*feed.Item
	logger ilog.Logger
}

func NewProvider(logger ilog.Logger, items ...*feed.Item) feed.Provider {
	if items == nil {
		items = []*feed.Item{}
	}
	return &provider{
		items:  items,
		logger: logger,
	}
}

func (p *provider) Type() string {
	return ProviderType
}

func (p *provider) Create(f *config.Feed) (feed.Service, error) {
	if f.Type == p.Type() {
		return NewService(f.Name, p.logger, p.items...), nil
	}
	return nil, nil
}
