package git_test

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	"github.com/patrickhuber/wrangle/internal/feed/git"
	"github.com/stretchr/testify/require"
)

func TestItemRepository(t *testing.T) {
	t.Run("List", func(t *testing.T) {
		repo := setupItemRepository(t)
		conformance.CanListAllItems(t, repo)
	})
	t.Run("Get", func(t *testing.T) {
		repo := setupItemRepository(t)
		conformance.CanGetItem(t, repo)
	})
}

func setupItemRepository(t *testing.T) feed.ItemRepository {
	workingDirectory := "/opt/wrangle/feed"
	fs := memfs.New()
	h := cross.NewTest(platform.Linux, arch.AMD64)
	path := h.Path()
	logger := log.Memory()
	repo := git.NewItemRepository(fs, path, logger, workingDirectory)
	items := conformance.GetItemList()
	for _, item := range items {
		require.Nil(t, repo.Save(item))
	}
	return repo
}
