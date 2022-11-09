package fs

import (
	"github.com/patrickhuber/go-log"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

func NewService(name string, fs filesystem.FileSystem, workingDirectory string, logger log.Logger) feed.Service {
	itemRepo := NewItemRepository(fs, workingDirectory)
	versionRepo := NewVersionRepository(fs, workingDirectory)
	return feed.NewService(name, itemRepo, versionRepo, logger)
}
