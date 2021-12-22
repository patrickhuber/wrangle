package memory_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("ItemRepository", func() {
	var (
		tester conformance.ItemRepositoryTester
	)
	BeforeEach(func() {
		items := map[string]*feed.Item{
			"test": {
				Package: &packages.Package{
					Name: "test",
					Versions: []*packages.PackageVersion{
						{
							Version: "1.0.0",
							Targets: []*packages.PackageTarget{
								{
									Platform:     "linux",
									Architecture: "amd64",
								},
							},
						},
					},
				},
			},
			"other": {
				Package: &packages.Package{
					Name: "other",
				},
			},
		}
		repo := memory.NewItemRepository(items)
		tester = conformance.NewItemRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all items", func() {
			tester.CanListAllItems(2)
		})
	})
	Describe("Get", func() {
		It("can get single item", func() {
			tester.CanGetPackage("other")
		})
	})
})
