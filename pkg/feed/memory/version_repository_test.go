package memory_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
)

func TestVersionRepository(t *testing.T) {
	t.Run("can get single version", func(t *testing.T) {
		tester := Setup(t)
		tester.CanGetSingleVersion(t)
	})
	t.Run("can list all versions", func(t *testing.T) {
		tester := Setup(t)
		tester.CanListAllVersions(t)
	})
	t.Run("can add a version", func(t *testing.T) {
		tester := Setup(t)
		tester.CanAddVersion(t, "test", "2.0.0")
	})
	t.Run("can update existing version", func(t *testing.T) {
		tester := Setup(t)
		tester.CanUpdateVersionNumber(t, "test", "1.0.0", "2.0.0")
	})
	t.Run("can add task", func(t *testing.T) {
		tester := Setup(t)
		tester.CanAddTask(t)
	})
}

func Setup(t *testing.T) conformance.VersionRepositoryTester {

	items := map[string]*feed.Item{}
	for _, i := range conformance.GetItemList() {
		items[i.Package.Name] = i
	}
	repo := memory.NewVersionRepository(items)
	return conformance.NewVersionRepositoryTester(repo)
}
