package feed

import (
	"strings"

	"github.com/spf13/afero"
)

type FeedServiceFactory interface {
	Get(packagesPath, feedURL string) FeedService
}

type feedServiceFactory struct {
	fs afero.Fs
}

func (factory *feedServiceFactory) Get(packagesPath, feedURL string) FeedService {
	if len(strings.TrimSpace(packagesPath)) == 0 {
		return NewGitFeedService(feedURL)
	}
	return NewFsFeedService(factory.fs, packagesPath)
}

func NewFeedServiceFactory(fs afero.Fs) FeedServiceFactory {
	return &feedServiceFactory{
		fs: fs,
	}
}
