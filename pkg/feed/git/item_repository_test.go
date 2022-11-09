package git_test

import (
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/git"
)

var _ = Describe("ItemRepository", func() {
	var (
		tester conformance.ItemRepositoryTester
	)
	BeforeEach(func() {
		workingDirectory := "/opt/wrangle/feed"
		fs := memfs.New()
		logger := log.Memory()
		repo := git.NewItemRepository(fs, logger, workingDirectory)
		items := conformance.GetItemList()
		for _, item := range items {
			Expect(repo.Save(item)).To(BeNil())
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
