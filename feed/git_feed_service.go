package feed

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
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
	response, err := svc.Get(&FeedGetRequest{Name: request.Name})
	if err != nil {
		return nil, err
	}
	latest := response.Package.Latest
	var latestPackageVersion *PackageVersion
	for _, packageVersion := range response.Package.Versions {
		if latest == packageVersion.Version {
			latestPackageVersion = packageVersion
		}
	}

	response.Package.Versions = []*PackageVersion{latestPackageVersion}
	return &FeedLatestResponse{
		Package: response.Package,
	}, nil
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
	latestVersions := map[string]*semver.Version{}
	pinLatestVersion := map[string]bool{}

	tree.Files().ForEach(func(f *object.File) error {

		segments := strings.Split(f.Name, "/")

		isLatestVersionFile := len(segments) == 3 && segments[0] == "feed" && segments[2] == "latest"
		isPackageVersionFile := len(segments) == 4 && segments[0] == "feed"
		if !isLatestVersionFile && !isPackageVersionFile {
			return nil
		}

		packageName := segments[1]

		if isLatestVersionFile {
			isLatestTagFile := segments[2] == "latest"
			if isLatestTagFile {
				// set the latest version in the array
				content, err := f.Contents()
				if err != nil {
					return err
				}
				// should this exit without error?
				ver, err := semver.NewVersion(content)
				if err != nil {
					return err
				}

				latestVersions[packageName] = ver
				pinLatestVersion[packageName] = true
				pkg, packageFound := packages[packageName]
				if packageFound {
					pkg.Latest = ver.String()
				}
			}
			return nil
		}

		if !isPackageVersionFile {
			return nil
		}

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

		ver, err := semver.NewVersion(packageVersion)
		if err != nil {
			return err
		}

		latest, latestFound := latestVersions[packageName]
		_, isPinned := pinLatestVersion[packageName]

		if !latestFound || (latest.Compare(ver) == -1 && !isPinned) {
			latestVersions[packageName] = ver
			latest = ver
			pkg.Latest = latest.String()
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
			version.Manifest.Content = content
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
