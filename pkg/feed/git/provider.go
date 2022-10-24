package git

import (
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

const ProviderType = "git"

type provider struct {
	logger ilog.Logger
}

func NewProvider(logger ilog.Logger) feed.Provider {
	return &provider{
		logger: logger,
	}
}

func (p *provider) Type() string {
	return ProviderType
}

func (p *provider) Create(f *config.Feed) (feed.Service, error) {
	if f.Type == p.Type() {
		return NewServiceFromURL(f.Name, f.URI, p.logger)
	}
	return nil, nil
}
