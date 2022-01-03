package fs_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	feedfs "github.com/patrickhuber/wrangle/pkg/feed/fs"
	"github.com/patrickhuber/wrangle/pkg/filesystem"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

var _ = Describe("VersionRespository", func() {
	var (
		tester conformance.VersionRepositoryTester
	)
	BeforeEach(func() {
		fs := filesystem.NewMemory()
		workingDirectory := "/opt/wrangle/feed"
		repo := feedfs.NewVersionRepository(fs, workingDirectory)
		versions := []*packages.Version{
			{
				Version: "1.0.0",
				Targets: []*packages.Target{
					{
						Platform:     "linux",
						Architecture: "amd64",
						Tasks: []*packages.Task{
							{
								Name: "download",
								Properties: map[string]string{
									"url": "https://www.google.com",
								},
							},
						},
					},
				},
			},
		}
		for _, version := range versions {
			err := repo.Save("test", version)
			Expect(err).To(BeNil())
		}
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
