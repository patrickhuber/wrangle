package fs

import (
	"github.com/patrickhuber/go-cross/filepath"
	crossfs "github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
)

type localServiceFactory struct {
	fs     crossfs.FS
	path   filepath.Provider
	logger log.Logger
}

func NewLocalServiceFactory(fs crossfs.FS, path filepath.Provider, logger log.Logger) feed.LocalServiceFactory {
	return &localServiceFactory{
		fs:     fs,
		path:   path,
		logger: logger,
	}
}

func (f *localServiceFactory) Create(directory string) (feed.Service, error) {
	feedDirectory := f.path.Join(directory, "feed")
	return NewService("local", f.fs, f.path, feedDirectory, f.logger), nil
}
