package git_test

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	gitfeed "github.com/patrickhuber/wrangle/internal/feed/git"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	setup := func(t *testing.T) conformance.ServiceTester {
		logger := log.Memory()
		store := memory.NewStorage()
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := memfs.New()
		path := h.Path()

		repository, err := git.Init(store, fs)
		require.NoError(t, err)

		svc, err := gitfeed.NewService("test", fs, repository, path, logger)
		require.NoError(t, err)

		items := conformance.GetItemList()
		response, err := svc.Update(&feed.UpdateRequest{
			Items: &feed.ItemUpdate{
				Add: items,
			},
		})
		require.NoError(t, err)
		require.Equal(t, conformance.TotalItemCount, response.Changed)
		return conformance.NewServiceTester(svc)
	}
	t.Run("can list all packages", func(t *testing.T) {
		tester := setup(t)
		tester.CanListAllPackages(t)
	})
	t.Run("can return latest version", func(t *testing.T) {
		tester := setup(t)
		tester.CanReturnLatestVersion(t)
	})
	t.Run("can return specific version", func(t *testing.T) {
		tester := setup(t)
		tester.CanReturnSpecificVersion(t)
	})
	t.Run("can add version", func(t *testing.T) {
		tester := setup(t)
		tester.CanAddVersion(t)
	})
	t.Run("can update existing version", func(t *testing.T) {
		tester := setup(t)
		tester.CanUpdateExistingVersion(t)
	})
}
