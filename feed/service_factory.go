package feed

import (
	"strings"

	"github.com/patrickhuber/wrangle/filesystem"
)

// ServiceFactory defines a contract for getting feed services
type ServiceFactory interface {
	Get(packagesPath, feedURL string) (Service, error)
}

type feedServiceFactory struct {
	fs filesystem.FileSystem
}

func (factory *feedServiceFactory) Get(packagesPath, feedURL string) (Service, error) {
	if len(strings.TrimSpace(packagesPath)) == 0 {
		svc, err := NewGitServiceFromURL(feedURL)
		if err != nil {
			return nil, err
		}
		return svc, err
	}
	return NewFsService(factory.fs, packagesPath), nil
}

// NewServiceFactory creates a new instance of the feed service factory with the given filesystem
func NewServiceFactory(fs filesystem.FileSystem) ServiceFactory {
	return &feedServiceFactory{
		fs: fs,
	}
}
