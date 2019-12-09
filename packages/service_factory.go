package packages

import (
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/ui"
)

// ServiceFactory defines a factory for creating Service instantces
type ServiceFactory interface {
	Get(feedService feed.FeedService) Service
}

type serviceFactory struct {
	console ui.Console
}

func (factory *serviceFactory) Get(feedService feed.FeedService) Service {
	return NewService(feedService, factory.console)
}

// NewServiceFactory creates a new package service factory
func NewServiceFactory(console ui.Console) ServiceFactory {
	return &serviceFactory{
		console: console,
	}
}
