package fs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
)

var _ = Describe("VersionRespository", func() {
	var (
		tester conformance.VersionRepositoryTester
	)
	BeforeEach(func() {
		fs := filesystem.NewMemory()
		workingDirectory := "/opt/wrangle/feed"
		Expect(fs.Write(crosspath.Join(workingDirectory, "test", "state.yml"), []byte(""), 0644)).To(BeNil())
		repo := feedfs.NewVersionRepository(fs, workingDirectory)
		tester = conformance.NewVersionRepositoryTester(repo)
	})
	Describe("List", func() {
		It("can list all versions", func() {
			tester.CanListAllVersions("test", 1)
		})
	})
	Describe("Get", func() {
		It("can get single version", func() {
			tester.CanGetSingleVersion("test", "1.0.0")
		})
	})
	Describe("Update", func() {
		It("can update ", func() {
			tester.CanUpdateVersionNumber("test", "1.0.0", "2.0.0")
		})
		It("can add task", func() {
			tester.CanAddTask("test", "1.0.0")
		})
		It("can", func() {
			tester.CanAddVersion("test", "2.0.0")
		})
	})
})
