package fs

import (
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

type itemRepository struct {
	fs               filesystem.FileSystem
	workingDirectory string
}

func NewItemRepository(fs filesystem.FileSystem, workingDirectory string) feed.ItemRepository {
	return &itemRepository{
		fs:               fs,
		workingDirectory: workingDirectory,
	}
}

func (r *itemRepository) List(include *feed.ItemGetInclude) ([]*feed.Item, error) {
	return nil, nil
}
func (r *itemRepository) Get(name string, include *feed.ItemGetInclude) (*feed.Item, error) {
	return nil, nil
}
func (r *itemRepository) Save(item *feed.Item, option *feed.ItemSaveOption) error {
	return nil
}
