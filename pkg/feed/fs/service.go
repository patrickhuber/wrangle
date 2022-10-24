package fs

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/ilog"
)

func NewService(name string, fs filesystem.FileSystem, workingDirectory string, logger ilog.Logger) feed.Service {
	itemRepo := NewItemRepository(fs, workingDirectory)
	versionRepo := NewVersionRepository(fs, workingDirectory)
	return feed.NewService(name, itemRepo, versionRepo, logger)
}
