package fs_test

import (
	. "github.com/onsi/ginkgo"

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
		tester = conformance.NewItemRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all items", func() {
			tester.CanListAllItems(1)
		})
	})
	Describe("Get", func() {
		It("can get package", func() {
			tester.CanGetPackage("test")
		})
	})
})
