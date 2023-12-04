package memory_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	"github.com/patrickhuber/wrangle/internal/feed/memory"
)

func TestItemRepository(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		ir := setupItemRepository()
		conformance.CanListAllItems(t, ir)
	})
	t.Run("get", func(t *testing.T) {
		ir := setupItemRepository()
		conformance.CanGetItem(t, ir)
	})
}

func setupItemRepository() feed.ItemRepository {
	items := map[string]*feed.Item{}
	for _, i := range conformance.GetItemList() {
		items[i.Package.Name] = i
	}
	repo := memory.NewItemRepository(items)
	return repo
}
