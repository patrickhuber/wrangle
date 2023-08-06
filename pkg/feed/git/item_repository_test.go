package git_test

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/host"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/git"
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
	h := host.NewTest(platform.Linux, arch.AMD64)
	path := h.Path
	logger := log.Memory()
	repo := git.NewItemRepository(fs, path, logger, workingDirectory)
	items := conformance.GetItemList()
	for _, item := range items {
		require.Nil(t, repo.Save(item))
	}
	return repo
}
