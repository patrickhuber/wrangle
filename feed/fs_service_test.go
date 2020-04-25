package feed_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/filesystem"

	"github.com/patrickhuber/wrangle/feed"
)

var _ = Describe("FsService", func() {
	var (
		fs          filesystem.FileSystem
		feedService feed.Service
		feedTest    FeedTest
		packages    []feed.Package
	)
	BeforeEach(func() {
		fs = filesystem.NewMemory()
		packages = []feed.Package{
			{
				Name: "bbr",
				Versions: []*feed.PackageVersion{
					{
						Version: "1.2.8",
					},
					{
						Version: "1.3.2",
						Manifest: &feed.PackageVersionManifest{
							Content: "Content",
							Name:    "Name",
						},
					},
				},
			},
			{
				Name:   "test",
				Latest: "1.0.1",
				Versions: []*feed.PackageVersion{
					{
						Version: "1.0.0",
					},
					{
						Version: "1.0.1",
					},
					{
						Version: "2.0.0",
					},
				},
			},
		}

		err := writePackagesToFileSystem(packages, fs)
		Expect(err).To(BeNil())

		feedService = feed.NewFsService(fs, "/wrangle/packages")
		feedTest = NewFeedTest(feedService)
	})
	Describe("List", func() {
		It("lists all packages", func() {
			feedTest.ListsExactPackages(packages)
		})
	})
	Describe("Get", func() {
		It("gets all versions by name", func() {
			feedTest.GetsAllVersionsByName("test", []string{
				"1.0.0",
				"1.0.1",
				"2.0.0",
			})
		})
		It("gets specific version by name and version", func() {
			feedTest.GetsSpecificVersionByNameAndVersion("test", "1.0.1")
		})
		Context("no package names match", func() {
			It("is empty", func() {
				feedTest.GetReturnsEmptyValueWhenNoPackageNameMatches("notFound")
			})
		})
		Context("no package no versions match", func() {
			It("version list is empty", func() {
				feedTest.GetReturnsEmptyValueWhenNoPackageVersionMatches("test", "2.0.1")
			})
		})
		It("returns content", func() {
			feedTest.GetReturnsContentWhenRequested("bbr", "1.3.2", "Content")
		})
	})
	Describe("Lastest", func() {
		It("gets latest package version", func() {
			feedTest.LastestReturnsLatestPackageVersion("bbr", "1.3.2")
		})
		Context("when tagged", func() {
			It("gets latest package version", func() {
				feedTest.LastestReturnsLatestPackageVersion("test", "1.0.1")
			})
		})
	})
	Describe("Info", func() {
		It("gets info for file system path", func() {
			feedTest.InfoReturnsURI("/wrangle/packages")
		})
	})
})

func writePackagesToFileSystem(packages []feed.Package, fs filesystem.FileSystem) error {
	// write packages to file system
	for p := 0; p < len(packages); p++ {
		pkg := packages[p]
		if pkg.Latest != "" {
			latestPath := fmt.Sprintf("/wrangle/packages/%s/latest", pkg.Name)
			if err := fs.Write(latestPath, []byte(pkg.Latest), 0666); err != nil {
				return err
			}
		}
		for v := 0; v < len(pkg.Versions); v++ {
			ver := pkg.Versions[v]
			path := fmt.Sprintf("/wrangle/packages/%s/%s/%s.%s.yml", pkg.Name, ver.Version, pkg.Name, ver.Version)
			content := ""
			if ver.Manifest != nil {
				content = ver.Manifest.Content
			}
			if err := fs.Write(path, []byte(content), 0666); err != nil {
				return err
			}
		}
	}
	return nil
}
