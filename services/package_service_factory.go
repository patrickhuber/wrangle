package services

import (
	"github.com/patrickhuber/wrangle/feed"
	"github.com/patrickhuber/wrangle/ui"
)

// PackageServiceFactory defines a factory for creating PackagesService instantces
type PackageServiceFactory interface {
	Get(feedService feed.FeedService) PackagesService
}

type packageServiceFactory struct {
	console ui.Console
}

func (factory *packageServiceFactory) Get(feedService feed.FeedService) PackagesService {
	return NewPackagesService(feedService, factory.console)
}

// NewPackageServiceFactory creates a new package service factory
func NewPackageServiceFactory(console ui.Console) PackageServiceFactory {
	return &packageServiceFactory{
		console: console,
	}
}
