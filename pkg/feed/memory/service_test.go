package memory_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("Service", func() {
	var (
		tester conformance.ServiceTester
	)
	BeforeEach(func() {
		service, err := memory.NewService(
			"test",
			&feed.Item{
				State: &feed.State{
					LatestVersion: "1.1.0",
				},
				Package: &packages.Package{
					Name: "test",
					Versions: []*packages.Version{
						{
							Version: "1.0.0",
						},
						{
							Version: "1.1.0",
						},
						{
							Version: "1.0.1",
						},
					},
				},
			})
		Expect(err).To(BeNil())
		tester = conformance.NewServiceTester(service)
	})
	Describe("List", func() {
		It("can list all packages", func() {
			tester.CanListAllPackages(1)
		})
		It("can return latest version", func() {
			tester.CanReturnLatestVersion("test", "1.1.0")
		})
		It("can return specific version", func() {
			tester.CanReturnSpecificVersion("test", "1.0.1")
		})
	})
	Describe("Update", func() {
		It("can add version", func() {
			tester.CanAddVersion("test", "2.0.0")
		})
		It("can update existing version", func() {
			tester.CanUpdateExistingVersion("test", "1.0.0", "2.0.0")
		})
	})
})
