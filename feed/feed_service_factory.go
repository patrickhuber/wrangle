package feed

import (
	"strings"

	"github.com/patrickhuber/wrangle/filesystem"
)

// FeedServiceFactory defines a contract for getting feed services
type FeedServiceFactory interface {
	Get(packagesPath, feedURL string) (FeedService, error)
}

type feedServiceFactory struct {
	fs filesystem.FileSystem
}

func (factory *feedServiceFactory) Get(packagesPath, feedURL string) (FeedService, error) {
	if len(strings.TrimSpace(packagesPath)) == 0 {
		svc, err := NewGitFeedServiceFromURL(feedURL)
		if err != nil {
			return nil, err
		}
		return svc, err
	}
	return NewFsFeedService(factory.fs, packagesPath), nil
}

// NewFeedServiceFactory creates a new instance of the feed service factory with the given filesystem
func NewFeedServiceFactory(fs filesystem.FileSystem) FeedServiceFactory {
	return &feedServiceFactory{
		fs: fs,
	}
}
