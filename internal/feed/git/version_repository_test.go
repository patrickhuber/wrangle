package git_test

import (
	"testing"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/go-xplat/setup"
	"github.com/patrickhuber/wrangle/internal/feed/conformance"
	"github.com/patrickhuber/wrangle/internal/feed/git"
	"github.com/stretchr/testify/require"
)

func TestVersionRepository(t *testing.T) {
	t.Run("list_all_versions", func(t *testing.T) {
		tester := CreateTest(t)
		tester.CanListAllVersions(t)
	})
	t.Run("get_single_version", func(t *testing.T) {
		tester := CreateTest(t)
		tester.CanGetSingleVersion(t)
	})
	t.Run("update_version", func(t *testing.T) {
		tester := CreateTest(t)
		tester.CanUpdateVersionNumber(t, "test", "1.0.0", "2.0.0")
	})
	t.Run("add_task", func(t *testing.T) {
		tester := CreateTest(t)
		tester.CanAddTask(t)
	})
	t.Run("add_version", func(t *testing.T) {
		tester := CreateTest(t)
		tester.CanAddVersion(t, "test", "2.0.0")
	})
}

func CreateTest(t *testing.T) conformance.VersionRepositoryTester {
	fs := memfs.New()
	workingDirectory := "/opt/wrangle/feed"
	logger := log.Memory()
	h := setup.NewTest(setup.Platform(platform.Linux))
	path := h.Path
	repo := git.NewVersionRepository(fs, logger, path, workingDirectory)
	items := conformance.GetItemList()
	for _, item := range items {
		for _, version := range item.Package.Versions {
			err := repo.Save(item.Package.Name, version)
			require.NoError(t, err)
		}
	}
	return conformance.NewVersionRepositoryTester(repo)
}
