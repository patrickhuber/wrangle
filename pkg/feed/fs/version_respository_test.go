package fs_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/fs"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
)

var _ = Describe("VersionRespository", func() {
	var (
		tester conformance.VersionRepositoryTester
	)
	BeforeEach(func() {
		fs := fs.NewMemory()
		path := filepath.NewProcessorWithPlatform(platform.Linux)
		workingDirectory := "/opt/wrangle/feed"
		repo := feedfs.NewVersionRepository(fs, path, workingDirectory)
		items := conformance.GetItemList()
		for _, item := range items {
			for _, version := range item.Package.Versions {
				err := repo.Save(item.Package.Name, version)
				Expect(err).To(BeNil())
			}
		}
		tester = conformance.NewVersionRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all versions", func() {
			tester.CanListAllVersions()
		})
	})
	Describe("Get", func() {
		It("can get single version", func() {
			tester.CanGetSingleVersion()
		})
	})
	Describe("Update", func() {
		It("can update ", func() {
			tester.CanUpdateVersionNumber("test", "1.0.0", "2.0.0")
		})
		It("can add task", func() {
			tester.CanAddTask()
		})
		It("can", func() {
			tester.CanAddVersion("test", "2.0.0")
		})
	})
})
