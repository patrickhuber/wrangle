package feed

import (
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-billy.v4/memfs"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

type gitFeedService struct {
	repository *git.Repository
	name       string
}

// NewGitFeedServiceFromURL returns a FeedService instance by cloning from the given URL
func NewGitFeedServiceFromURL(URL string) (FeedService, error) {
	fs := memfs.New()
	storer := memory.NewStorage()

	repository, err := git.Clone(storer, fs, &git.CloneOptions{
		URL: URL,
	})
	if err != nil {
		return nil, err
	}

	return NewGitFeedService(repository), nil
}

// NewGitFeedService creates a new FeedService from the given repositry
func NewGitFeedService(repository *git.Repository) FeedService {
	return &gitFeedService{
		repository: repository,
		name:       "remote",
	}
}

func (svc *gitFeedService) List(request *FeedListRequest) (*FeedListResponse, error) {
	packages, err := svc.find(nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeedListResponse{
		Packages: packages,
	}, nil
}

func (svc *gitFeedService) Get(request *FeedGetRequest) (*FeedGetResponse, error) {
	versions := []string{}
	if strings.TrimSpace(request.Version) != "" {
		versions = append(versions, request.Version)
	}
	where := &packageCriteriaWhere{
		Or: []*packageCriteriaAnd{
			&packageCriteriaAnd{
				And: []*packageCriteria{
					&packageCriteria{
						Name:     request.Name,
						Versions: versions,
					},
				},
			},
		},
	}
	include := &packageInclude{
		Content: request.IncludeContent,
	}
	packages, err := svc.find(where, include)
	if err != nil {
		return nil, err
	}

	var pkg *Package
	if len(packages) > 0 {
		pkg = packages[0]
	}

	return &FeedGetResponse{
		Package: pkg,
	}, nil

}

func (svc *gitFeedService) Create(request *FeedCreateRequest) (*FeedCreateResponse, error) {
	return nil, nil
}

func (svc *gitFeedService) Latest(request *FeedLatestRequest) (*FeedLatestResponse, error) {
	return nil, nil
}

func (svc *gitFeedService) find(where *packageCriteriaWhere, include *packageInclude) ([]*Package, error) {

	ref, err := svc.repository.Head()
	if err != nil {
		return nil, err
	}

	commit, err := svc.repository.CommitObject(ref.Hash())
	if err != nil {
		return nil, err
	}

	tree, err := commit.Tree()
	if err != nil {
		return nil, err
	}

	packages := map[string]*Package{}
	tree.Files().ForEach(func(f *object.File) error {

		segments := strings.Split(f.Name, "/")

		if len(segments) != 4 {
			return nil
		}

		if segments[0] != "feed" {
			return nil
		}

		packageName := segments[1]
		packageVersion := segments[2]
		packageVersionManifestFile := segments[3]

		packageVersionManifestName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)

		if packageVersionManifestName != packageVersionManifestFile {
			return nil
		}

		if !evaluate(where, packageName, packageVersion) {
			return nil
		}

		pkg, ok := packages[packageName]
		if !ok {
			pkg = &Package{
				Name:     packageName,
				Versions: []*PackageVersion{},
			}
			packages[packageName] = pkg
		}

		version := &PackageVersion{
			Manifest: &PackageVersionManifest{
				Name: packageVersionManifestFile},
			Version: packageVersion,
			Feeds:   []string{svc.name},
		}
		if include != nil && include.Content {
			content, err := f.Contents()
			if err != nil {
				return err
			}
			version.Manifest.Content = &content
		}
		pkg.Versions = append(pkg.Versions, version)

		return nil
	})

	pkgList := []*Package{}
	for _, pkg := range packages {
		pkgList = append(pkgList, pkg)
	}
	return pkgList, nil
}
