package fs_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
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

	t.Run("reads numeric architecture values as strings", func(t *testing.T) {
		workingDirectory := "/opt/wrangle/feed"
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		itemPath := path.Join(workingDirectory, "task")
		require.NoError(t, fs.MkdirAll(itemPath, 0755))
		require.NoError(t, fs.WriteFile(path.Join(itemPath, "platforms.yml"), []byte(`platforms:
- platform: linux
  architectures: [386, amd64]
`), 0644))

		repo := feedfs.NewItemRepository(fs, path, workingDirectory, log.Memory())
		item, err := repo.Get("task",
			feed.ItemGetState(false),
			feed.ItemGetTemplate(false),
			feed.ItemGetResource(false),
		)
		require.NoError(t, err)
		require.Equal(t, []string{"386", "amd64"}, item.Platforms[0].Architectures)
	})

	t.Run("skips directories without resource files", func(t *testing.T) {
		workingDirectory := "/opt/wrangle/feed"
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		require.NoError(t, fs.MkdirAll(path.Join(workingDirectory, "linuxkit", "0.6"), 0755))

		var output bytes.Buffer
		logger := log.MemoryWith(&output, log.WithLevel(log.TraceLevel))
		repo := feedfs.NewItemRepository(fs, path, workingDirectory, logger)
		items, err := repo.List(feed.ItemGetSkipMissingResource())
		require.NoError(t, err)
		require.Empty(t, items)
		require.Contains(t, output.String(), "skipping feed package linuxkit")
	})
}

func setupItemRepository(t *testing.T) feed.ItemRepository {
	workingDirectory := "/opt/wrangle/feed"
	h := cross.NewTest(platform.Linux, arch.AMD64)
	fs := h.FS()
	path := h.Path()
	repo := feedfs.NewItemRepository(fs, path, workingDirectory, log.Memory())

	items := conformance.GetItemList()
	for _, item := range items {
		err := repo.Save(item)
		require.NoError(t, err)
	}
	return repo
}
