package feed

import (
	"strings"

	"github.com/patrickhuber/wrangle/filesystem"
)

type FeedServiceFactory interface {
	Get(packagesPath, feedURL string) FeedService
}

type feedServiceFactory struct {
	fs filesystem.FileSystem
}

func (factory *feedServiceFactory) Get(packagesPath, feedURL string) FeedService {
	if len(strings.TrimSpace(packagesPath)) == 0 {
		return NewGitFeedService(feedURL)
	}
	return NewFsFeedService(factory.fs, packagesPath)
}

func NewFeedServiceFactory(fs filesystem.FileSystem) FeedServiceFactory {
	return &feedServiceFactory{
		fs: fs,
	}
}
