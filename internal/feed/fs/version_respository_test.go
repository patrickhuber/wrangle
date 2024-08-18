package fs_test

import (
	"testing"

	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/internal/feed/fs"
	"github.com/stretchr/testify/require"
)

func TestVersionRepository(t *testing.T) {
	t.Run("can list all versions", func(t *testing.T) {
		tester := Setup(t)
		tester.CanListAllVersions(t)
	})
	t.Run("can get single version", func(t *testing.T) {
		tester := Setup(t)
		tester.CanGetSingleVersion(t)
	})
	t.Run("can update ", func(t *testing.T) {
		tester := Setup(t)
		tester.CanUpdateVersionNumber(t, "test", "1.0.0", "2.0.0")
	})
	t.Run("can add task", func(t *testing.T) {
		tester := Setup(t)
		tester.CanAddTask(t)
	})
	t.Run("can", func(t *testing.T) {
		tester := Setup(t)
		tester.CanAddVersion(t, "test", "2.0.0")
	})
}

func Setup(t *testing.T) conformance.VersionRepositoryTester {
	h := setup.NewTest(setup.Platform(platform.Linux))
	fs := h.FS
	path := h.Path
	workingDirectory := "/opt/wrangle/feed"
	repo := feedfs.NewVersionRepository(fs, path, workingDirectory)
	items := conformance.GetItemList()
	for _, item := range items {
		for _, version := range item.Package.Versions {
			err := repo.Save(item.Package.Name, version)
			require.NoError(t, err)
		}
	}
	tester := conformance.NewVersionRepositoryTester(repo)
	return tester
}
