package memory_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
)

var _ = Describe("PackageVersionRepository", func() {
	var (
		tester conformance.VersionRepositoryTester
	)
	BeforeEach(func() {
		items := map[string]*feed.Item{}
		for _, i := range conformance.GetItemList() {
			items[i.Package.Name] = i
		}
		repo := memory.NewVersionRepository(items)
		tester = conformance.NewVersionRepositoryTester(repo)
	})
	Describe("Get", func() {
		It("can get single version", func() {
			tester.CanGetSingleVersion("test", "1.0.0")
		})
	})
	Describe("List", func() {
		It("can list all versions", func() {
			tester.CanListAllVersions("test", 3)
		})
	})
	Describe("Update", func() {
		It("can add a version", func() {
			tester.CanAddVersion("test", "2.0.0")
		})
		It("can update existing version", func() {
			tester.CanUpdateVersionNumber("test", "1.0.0", "2.0.0")
		})
		It("can add target", func() {
			tester.CanAddTarget("test", "1.0.0")
		})
		It("can add task", func() {
			tester.CanAddTask("test", "1.0.0")
		})
	})
})
