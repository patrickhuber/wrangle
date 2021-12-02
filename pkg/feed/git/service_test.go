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
	gitfeed "github.com/patrickhuber/wrangle/pkg/feed/git"
	"github.com/patrickhuber/wrangle/pkg/packages"
	"gopkg.in/yaml.v2"
)

var _ = Describe("GitService", func() {
	Describe("List", func() {
		It("can list all packages", func() {
			store := memory.NewStorage()
			fs := memfs.New()

			repository, err := git.Init(store, fs)
			Expect(err).To(BeNil())

			items := []feed.Item{
				{
					Package: &packages.Package{
						Name: "test",
						Versions: []*packages.PackageVersion{
							{
								Version: "1.0.0",
							},
						},
					},
					State: &feed.State{
						LatestVersion: "1.0.0",
					},
					Template: "",
					Platforms: []*feed.Platform{
						{
							Name:          "windows",
							Architectures: []string{"amd64", "386"},
						},
					},
				},
				{
					Package: &packages.Package{
						Name: "ffa",
						Versions: []*packages.PackageVersion{
							{
								Version: "1.0.0",
							},
						},
					},
					State: &feed.State{
						LatestVersion: "1.0.0",
					},
					Template: "",
					Platforms: []*feed.Platform{
						{
							Name:          "windows",
							Architectures: []string{"amd64", "386"},
						},
					},
				},
				{
					Package: &packages.Package{
						Name: "tsa",
						Versions: []*packages.PackageVersion{
							{
								Version: "1.0.0",
							},
						},
					},
					State: &feed.State{
						LatestVersion: "1.0.0",
					},
					Template: "",
					Platforms: []*feed.Platform{
						{
							Name:          "windows",
							Architectures: []string{"amd64", "386"},
						},
					},
				},
			}
			err = writeItemsToGitRepo(items, repository)
			Expect(err).To(BeNil())

			svc, err := gitfeed.NewService(repository)
			Expect(err).To(BeNil())

			response, err := svc.List(&feed.ListRequest{})
			Expect(err).To(BeNil())
			Expect(response).ToNot(BeNil())
			Expect(len(response.Items)).To(Equal(3))
		})
	})
})

type Manifest struct {
	Package *ManifestPackage
}

type ManifestPackage struct {
	Name    string
	Version string
	Targets []*packages.PackageTarget
}

func writeItemsToGitRepo(items []feed.Item, repository *git.Repository) error {
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
			manifest := &Manifest{
				Package: &ManifestPackage{
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
