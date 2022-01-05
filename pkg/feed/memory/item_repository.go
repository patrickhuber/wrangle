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

func (r *itemRepository) Get(name string, options ...feed.ItemGetOption) (*feed.Item, error) {
	item, ok := r.items[name]
	if !ok {
		return nil, nil
	}
	return item, nil
}

func (r *itemRepository) List(options ...feed.ItemGetOption) ([]*feed.Item, error) {
	items := []*feed.Item{}
	for _, item := range r.items {
		items = append(items, item)
	}
	return items, nil
}

func (r *itemRepository) Save(item *feed.Item, option ...feed.ItemSaveOption) error {
	r.items[item.Package.Name] = item
	return nil
}
