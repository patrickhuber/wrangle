package memory

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
)

func NewService(name string, items ...*feed.Item) (feed.Service, error) {
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

	return feed.NewService(name, itemRepo, versionRepo), nil
}
