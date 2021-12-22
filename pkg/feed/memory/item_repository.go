package memory

import "github.com/patrickhuber/wrangle/pkg/feed"

type itemRepository struct {
	items map[string]*feed.Item
}

func NewItemRepository(items map[string]*feed.Item) feed.ItemRepository {
	return &itemRepository{
		items: items,
	}
}

func (r *itemRepository) Get(name string) (*feed.Item, error) {
	item, ok := r.items[name]
	if !ok {
		return nil, nil
	}
	return item, nil
}

func (r *itemRepository) List(where []*feed.ItemReadAnyOf) ([]*feed.Item, error) {
	items := []*feed.Item{}
	for name, item := range r.items {
		if feed.IsMatch(where, name) {
			items = append(items, item)
		}
	}
	return items, nil
}
