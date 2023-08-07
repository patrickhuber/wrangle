package fs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
)

func TestItemRepo(t *testing.T) {

	t.Run("list", func(t *testing.T) {
		ir := setupItemRepository(t)
		conformance.CanListAllItems(t, ir)
	})

	t.Run("get", func(t *testing.T) {
		ir := setupItemRepository(t)
		conformance.CanGetItem(t, ir)
	})
}

func setupItemRepository(t *testing.T) feed.ItemRepository {
	workingDirectory := "/opt/wrangle/feed"
	h := setup.NewTest(setup.Platform(platform.Linux))
	fs := h.FS
	path := h.Path
	repo := feedfs.NewItemRepository(fs, path, workingDirectory)

	items := conformance.GetItemList()
	for _, item := range items {
		err := repo.Save(item)
		require.Nil(t, err)
	}
	return repo
}
