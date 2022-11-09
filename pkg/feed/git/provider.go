package git

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/pkg/config"
	"github.com/patrickhuber/wrangle/pkg/feed"
)

const ProviderType = "git"

type provider struct {
	logger log.Logger
}

func NewProvider(logger log.Logger) feed.Provider {
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
