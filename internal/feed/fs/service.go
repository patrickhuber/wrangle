package fs

import (
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-log"

	"github.com/patrickhuber/wrangle/internal/feed"
)

func NewService(name string, fs fs.FS, path filepath.Provider, workingDirectory string, logger log.Logger) feed.Service {
	itemRepo := NewItemRepository(fs, path, workingDirectory)
	versionRepo := NewVersionRepository(fs, path, workingDirectory)
	return feed.NewService(name, itemRepo, versionRepo, logger)
}
