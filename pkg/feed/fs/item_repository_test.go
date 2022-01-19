package fs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

var _ = Describe("ItemRepository", func() {
	var (
		tester conformance.ItemRepositoryTester
	)
	BeforeEach(func() {
		workingDirectory := "/opt/wrangle/feed"
		fs := filesystem.NewMemory()
		repo := feedfs.NewItemRepository(fs, workingDirectory)

		items := conformance.GetItemList()
		for _, item := range items {
			err := repo.Save(item)
			Expect(err).To(BeNil())
		}
		tester = conformance.NewItemRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all items", func() {
			tester.CanListAllItems()
		})
	})
	Describe("Get", func() {
		It("can get package", func() {
			tester.CanGetItem()
		})
	})
})
