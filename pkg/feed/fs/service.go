package fs

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"

	"github.com/patrickhuber/wrangle/pkg/feed"
)

func NewService(name string, fs fs.FS, path *filepath.Processor, workingDirectory string, logger log.Logger) feed.Service {
	itemRepo := NewItemRepository(fs, path, workingDirectory)
	versionRepo := NewVersionRepository(fs, path, workingDirectory)
	return feed.NewService(name, itemRepo, versionRepo, logger)
}
