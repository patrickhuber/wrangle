package git

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/feed"
)

const ProviderType = "git"

type provider struct {
	logger log.Logger
	path   *filepath.Processor
}

func NewProvider(path *filepath.Processor, logger log.Logger) feed.Provider {
	return &provider{
		logger: logger,
		path:   path,
	}
}

func (p *provider) Type() string {
	return ProviderType
}

func (p *provider) Create(f config.Feed) (feed.Service, error) {
	if f.Type == p.Type() {
		return NewServiceFromURL(f.Name, f.URI, p.path, p.logger)
	}
	return nil, nil
}
