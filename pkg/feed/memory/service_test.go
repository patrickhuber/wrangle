package memory_test

import (
	. "github.com/onsi/ginkgo/v2"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	"github.com/patrickhuber/wrangle/pkg/feed/memory"
)

var _ = Describe("Service", func() {
	var (
		tester conformance.ServiceTester
	)
	BeforeEach(func() {
		items := conformance.GetItemList()
		service := memory.NewService("test", items...)
		tester = conformance.NewServiceTester(service)
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
