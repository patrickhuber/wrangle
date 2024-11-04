package fs_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/internal/feed/fs"
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
	h := cross.NewTest(platform.Linux, arch.AMD64)
	fs := h.FS()
	path := h.Path()
	repo := feedfs.NewItemRepository(fs, path, workingDirectory)

	items := conformance.GetItemList()
	for _, item := range items {
		err := repo.Save(item)
		require.Nil(t, err)
	}
	return repo
}
