package git_test

import (
	"fmt"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/util"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/patrickhuber/wrangle/pkg/crosspath"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/feed/conformance"
	gitfeed "github.com/patrickhuber/wrangle/pkg/feed/git"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
)

var _ = Describe("GitService", func() {
	var (
		tester conformance.ServiceTester
	)
	BeforeEach(func() {
		store := memory.NewStorage()
		fs := memfs.New()

		repository, err := git.Init(store, fs)
		Expect(err).To(BeNil())
		items := conformance.GetItemList()
		err = writeItemsToGitRepo(items, repository)
		Expect(err).To(BeNil())

		svc, err := gitfeed.NewService("test", fs, repository)
		Expect(err).To(BeNil())
		tester = conformance.NewServiceTester(svc)
	})
	Describe("List", func() {
		It("can list all packages", func() {
			tester.CanListAllPackages(3)
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

func writeItemsToGitRepo(items []*feed.Item, repository *git.Repository) error {
	worktree, err := repository.Worktree()
	if err != nil {
		return err
	}
	fs := worktree.Filesystem

	for _, i := range items {
		packagePath := crosspath.Join("feed", i.Package.Name)
		err = fs.MkdirAll(packagePath, 0600)
		if err != nil {
			return err
		}
		for _, v := range i.Package.Versions {
			manifest := &packages.Manifest{
				Package: &packages.ManifestPackage{
					Name:    i.Package.Name,
					Version: v.Version,
					Targets: v.Targets,
				},
			}

			content, err := yaml.Marshal(manifest)
			if err != nil {
				return err
			}

			versionPath := crosspath.Join(packagePath, manifest.Package.Version)
			err = fs.MkdirAll(versionPath, 0600)
			if err != nil {
				return err
			}

			filePath := fmt.Sprintf("%s/package.yml", versionPath)
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
