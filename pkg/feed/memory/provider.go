package memory

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

const ProviderType = "memory"

type provider struct {
	items []*feed.Item
}

func NewProvider(items ...*feed.Item) feed.Provider {
	return &provider{
		items: items,
	}
}

func (p *provider) Type() string {
	return ProviderType
}

func (p *provider) Create(f *config.Feed) (feed.Service, error) {
	if f.Type == p.Type() {
		return NewService(f.Name, p.items...)
	}
	return nil, nil
}
