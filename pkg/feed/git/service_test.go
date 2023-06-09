package git_test

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/go-xplat/filepath"
	"github.com/patrickhuber/go-xplat/platform"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	gitfeed "github.com/patrickhuber/wrangle/pkg/feed/git"
)

var _ = Describe("GitService", func() {
	var (
		tester conformance.ServiceTester
	)
	BeforeEach(func() {
		logger := log.Memory()
		store := memory.NewStorage()
		fs := memfs.New()
		path := filepath.NewProcessorWithPlatform(platform.Linux)

		repository, err := git.Init(store, fs)
		Expect(err).To(BeNil())

		svc, err := gitfeed.NewService("test", fs, repository, path, logger)
		Expect(err).To(BeNil())

		items := conformance.GetItemList()
		response, err := svc.Update(&feed.UpdateRequest{
			Items: &feed.ItemUpdate{
				Add: items,
			},
		})
		Expect(err).To(BeNil())
		Expect(response.Changed).To(Equal(conformance.TotalItemCount))
		tester = conformance.NewServiceTester(svc)
	})
	Describe("List", func() {
		It("can list all packages", func() {
			tester.CanListAllPackages()
		})
		It("can return latest version", func() {
			tester.CanReturnLatestVersion()
		})
		It("can return specific version", func() {
			tester.CanReturnSpecificVersion()
		})
	})
	Describe("Update", func() {
		It("can add version", func() {
			tester.CanAddVersion()
		})
		It("can update existing version", func() {
			tester.CanUpdateExistingVersion()
		})
	})
})
