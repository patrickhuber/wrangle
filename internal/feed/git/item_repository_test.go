package git_test

import (
	"io/fs"
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
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

func TestItemRepositoryMissingMetadata(t *testing.T) {
	for _, missing := range []string{"platforms.yml", "state.yml", "template.yml"} {
		t.Run(missing, func(t *testing.T) {
			filesystem := memfs.New()
			h := cross.NewTest(platform.Linux, arch.AMD64)
			repo := git.NewItemRepository(filesystem, h.Path(), log.Memory(), "/feed")
			require.NoError(t, filesystem.MkdirAll("/feed/example", 0755))
			for _, name := range []string{"platforms.yml", "state.yml", "template.yml"} {
				if name != missing {
					require.NoError(t, util.WriteFile(filesystem, "/feed/example/"+name, []byte("{}"), 0644))
				}
			}

			_, err := repo.List()
			require.ErrorIs(t, err, fs.ErrNotExist)
			require.ErrorContains(t, err, "/feed/example/"+missing)
		})
	}
}

func TestItemRepositoryListAllowsMissingState(t *testing.T) {
	filesystem := memfs.New()
	h := cross.NewTest(platform.Linux, arch.AMD64)
	repo := git.NewItemRepository(filesystem, h.Path(), log.Memory(), "/feed")
	require.NoError(t, filesystem.MkdirAll("/feed/example", 0755))
	for _, name := range []string{"platforms.yml", "template.yml"} {
		require.NoError(t, util.WriteFile(filesystem, "/feed/example/"+name, []byte("{}"), 0644))
	}

	items, err := repo.List(feed.ItemGetStateOptional())
	require.NoError(t, err)
	require.Nil(t, items[0].State)
}

func TestItemRepositoryMissingFeedDirectory(t *testing.T) {
	h := cross.NewTest(platform.Linux, arch.AMD64)
	repo := git.NewItemRepository(memfs.New(), h.Path(), log.Memory(), "/missing-feed")

	_, err := repo.List()
	require.ErrorIs(t, err, fs.ErrNotExist)
	require.ErrorContains(t, err, "/missing-feed")
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
