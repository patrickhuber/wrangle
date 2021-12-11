package git

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

const ProviderType = "git"

type provider struct {
}

func NewProvider() feed.Provider {
	return &provider{}
}

func (p *provider) Type() string {
	return ProviderType
}

func (p *provider) Create(f *config.Feed) (feed.Service, error) {
	if f.Type == p.Type() {
		return NewServiceFromURL(f.Name, f.URI)
	}
	return nil, nil
}
