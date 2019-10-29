package feed_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/patrickhuber/wrangle/filesystem"

	"github.com/patrickhuber/wrangle/feed"
)

var _ = Describe("FeedService", func() {
	var (
		fs          filesystem.FileSystem
		feedService feed.FeedService
		feedTest    FeedTest
		packages    []feed.Package
	)
	BeforeEach(func() {
		fs = filesystem.NewMemory()
		packages = []feed.Package{
			feed.Package{
				Name: "test",
				Versions: []*feed.PackageVersion{
					&feed.PackageVersion{
						Version: "1.0.0",
					},
					&feed.PackageVersion{
						Version: "1.0.1",
					},
				},
			},
			feed.Package{
				Name: "other",
				Versions: []*feed.PackageVersion{
					&feed.PackageVersion{
						Version: "1.0.0",
						Manifest: &feed.PackageVersionManifest{
							Name:    "Name",
							Content: "Content",
						},
					},
				},
			},
			feed.Package{
				Name: "last",
				Versions: []*feed.PackageVersion{
					&feed.PackageVersion{
						Version: "1.0.0",
					},
				},
			},
		}

		err := writePackagesToFileSystem(packages, fs)
		Expect(err).To(BeNil())

		feedService = feed.NewFsFeedService(fs, "/wrangle/packages")
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
				feedTest.GetReturnsEmptyValueWhenNoPackageVersionMatches("test", "2.0.0")
			})
		})
		It("returns content", func() {
			feedTest.GetReturnsContentWhenRequested("other", "1.0.0", "Content")
		})
	})
	Describe("Lastest", func() {
		It("gets latest package version", func() {
			feedTest.LastestReturnsLatestPackageVersion("test", "1.0.1")
		})
	})
})

func writePackagesToFileSystem(packages []feed.Package, fs filesystem.FileSystem) error {
	// write packages to file system
	for p := 0; p < len(packages); p++ {
		pkg := packages[p]
		for v := 0; v < len(pkg.Versions); v++ {
			ver := pkg.Versions[v]
			path := fmt.Sprintf("/wrangle/packages/%s/%s/%s.%s.yml", pkg.Name, ver.Version, pkg.Name, ver.Version)
			content := ""
			if ver.Manifest != nil {
				content = ver.Manifest.Content
			}
			err := fs.Write(path, []byte(content), 0666)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
