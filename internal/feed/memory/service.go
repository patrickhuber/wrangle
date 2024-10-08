package memory

import (
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
)

func NewService(name string, logger log.Logger, items ...*feed.Item) feed.Service {
	itemMap := map[string]*feed.Item{}
	for _, i := range items {
		if i == nil || i.Package == nil || i.Package.Name == "" {
			continue
		}
		itemMap[i.Package.Name] = i
	}
	itemRepo := &itemRepository{
		items: itemMap,
	}
	versionRepo := &versionRepository{
		items: itemMap,
	}

	return feed.NewService(name, itemRepo, versionRepo, logger)
}
