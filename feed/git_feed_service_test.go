package feed_test

import (
	"fmt"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/patrickhuber/wrangle/feed"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-billy.v4/util"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var _ = Describe("GitFeedService", func() {
	var (
		svc    feed.FeedService
		tester FeedTest
	)
	BeforeEach(func() {
		store := memory.NewStorage()
		fs := memfs.New()

		repository, err := git.Init(store, fs)
		Expect(err).To(BeNil())

		packages := []feed.Package{
			feed.Package{
				Name: "bbr",
				Versions: []*feed.PackageVersion{
					&feed.PackageVersion{
						Version: "1.2.8",
					},
					&feed.PackageVersion{
						Version: "1.3.2",
						Manifest: &feed.PackageVersionManifest{
							Content: "Content",
							Name:    "Name",
						},
					},
				},
			},
			feed.Package{
				Name:   "test",
				Latest: "1.0.1",
				Versions: []*feed.PackageVersion{
					&feed.PackageVersion{
						Version: "1.0.0",
					},
					&feed.PackageVersion{
						Version: "1.0.1",
					},
					&feed.PackageVersion{
						Version: "2.0.0",
					},
				},
			},
		}
		writePackagesToGitRepository(packages, repository)

		svc = feed.NewGitFeedService(repository)
		Expect(svc).ToNot(BeNil())
		tester = NewFeedTest(svc)
	})
	Describe("List", func() {
		It("lists all packages", func() {
			tester.ListsAllPackages(2)
		})
	})

	Describe("Get", func() {
		It("gets all versions by name", func() {
			tester.GetsAllVersionsByName("bbr", []string{"1.2.8", "1.3.2"})
		})
		It("gets specific version by name and version", func() {
			tester.GetsSpecificVersionByNameAndVersion("bbr", "1.2.8")
		})
		Context("no package names match", func() {
			It("is empty", func() {
				tester.GetReturnsEmptyValueWhenNoPackageNameMatches("notFoundPackage")
			})
		})
		Context("no package no versions match", func() {
			It("is empty", func() {
				tester.GetReturnsEmptyValueWhenNoPackageVersionMatches("bbr", "2.0.0")
			})
		})
		It("returns content", func() {
			tester.GetReturnsContentWhenRequested("bbr", "1.3.2", "Content")
		})
	})
	Describe("Lastest", func() {
		It("gets latest package version", func() {
			tester.LastestReturnsLatestPackageVersion("bbr", "1.3.2")
		})
		Context("when tagged", func() {
			It("gets latest package version", func() {
				tester.LastestReturnsLatestPackageVersion("test", "1.0.1")
			})
		})
	})
})

func writePackagesToGitRepository(packages []feed.Package, repository *git.Repository) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	fs := worktree.Filesystem

	for p := 0; p < len(packages); p++ {
		pkg := packages[p]

		packagePath := fs.Join("feed", pkg.Name)

		err = fs.MkdirAll(packagePath, 0600)
		if err != nil {
			return nil
		}

		if pkg.Latest != "" {
			latestFilePath := fs.Join(packagePath, "latest")
			err = util.WriteFile(fs, latestFilePath, []byte(pkg.Latest), 0644)
			if err != nil {
				return err
			}
			_, err := worktree.Add(latestFilePath)
			if err != nil {
				return err
			}
		}

		for v := 0; v < len(pkg.Versions); v++ {
			ver := pkg.Versions[v]
			fileName := fmt.Sprintf("%s.%s.yml", pkg.Name, ver.Version)
			directoryPath := fs.Join(packagePath, ver.Version)

			err = fs.MkdirAll(directoryPath, 0600)
			if err != nil {
				return err
			}

			filePath := fs.Join(directoryPath, fileName)

			content := []byte("")
			if ver.Manifest != nil && ver.Manifest.Content != "" {
				content = []byte(ver.Manifest.Content)
			}

			err = util.WriteFile(fs, filePath, content, 0644)
			if err != nil {
				return err
			}

			_, err = worktree.Add(filePath)
			if err != nil {
				return err
			}
		}
	}

	_, err = worktree.Commit("initial revision", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "John Doe",
			Email: "john@doe.org",
			When:  time.Now(),
		},
	})

	return err
}
