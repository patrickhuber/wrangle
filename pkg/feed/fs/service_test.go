package fs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"

	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

var _ = Describe("Service", func() {
	var (
		tester conformance.ServiceTester
	)
	BeforeEach(func() {
		fs := filesystem.NewMemory()
		svc := feedfs.NewService("test", fs, "/opt/wrangle/feed")
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
