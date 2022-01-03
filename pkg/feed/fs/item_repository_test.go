package fs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("ItemRepository", func() {
	var (
		tester conformance.ItemRepositoryTester
	)
	BeforeEach(func() {
		workingDirectory := "/opt/wrangle/feed"
		fs := filesystem.NewMemory()
		repo := feedfs.NewItemRepository(fs, workingDirectory)
		option := &feed.ItemSaveOption{
			Platforms: true,
			State:     true,
			Template:  true,
		}
		items := []*feed.Item{
			{
				Package: &packages.Package{
					Name: "test",
					Versions: []*packages.Version{
						{
							Version: "1.0.0",
						},
					},
				},
				State: &feed.State{
					LatestVersion: "1.0.0",
				},
				Template: "",
				Platforms: []*feed.Platform{
					{
						Name:          "linux",
						Architectures: []string{"amd64"},
					},
				},
			},
		}
		for _, item := range items {
			err := repo.Save(item, option)
			Expect(err).To(BeNil())
		}
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
