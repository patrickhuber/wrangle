package memory_test

import (
	. "github.com/onsi/ginkgo"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
)

var _ = Describe("ItemRepository", func() {
	var (
		tester conformance.ItemRepositoryTester
	)
	BeforeEach(func() {
		items := map[string]*feed.Item{}
		for _, i := range conformance.GetItemList() {
			items[i.Package.Name] = i
		}
		repo := memory.NewItemRepository(items)
		tester = conformance.NewItemRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all items", func() {
			tester.CanListAllItems()
		})
	})
	Describe("Get", func() {
		It("can get single item", func() {
			tester.CanGetItem()
		})
	})
})
