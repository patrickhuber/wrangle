package feed

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/filepath"

	"github.com/patrickhuber/wrangle/filesystem"
)

type fsFeedService struct {
	fs   filesystem.FileSystem
	path string
	name string
}

// NewFsFeedService defines a feed service over the filesystem
func NewFsFeedService(fs filesystem.FileSystem, path string) FeedService {
	return &fsFeedService{
		fs:   fs,
		path: path,
		name: "local",
	}
}

func (svc *fsFeedService) List(request *FeedListRequest) (*FeedListResponse, error) {
	packages, err := svc.find(nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeedListResponse{
		Packages: packages,
	}, nil
}

func (svc *fsFeedService) Get(request *FeedGetRequest) (*FeedGetResponse, error) {
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

func (svc *fsFeedService) find(where *packageCriteriaWhere, include *packageInclude) ([]*Package, error) {
	packageFolders, err := svc.fs.ReadDir(svc.path)
	if err != nil {
		return nil, err
	}
	packages := []*Package{}
	for _, packageFolder := range packageFolders {
		if !packageFolder.IsDir() {
			continue
		}
		packageName := packageFolder.Name()
		packagePath := filepath.Join(svc.path, packageFolder.Name())
		packageVersions, err := svc.fs.ReadDir(packagePath)
		if err != nil {
			return nil, err
		}

		versions := []*PackageVersion{}
		for _, packageVersionFolder := range packageVersions {
			if !packageVersionFolder.IsDir() {
				continue
			}
			packageVersion := packageVersionFolder.Name()
			packageVersionManifestName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)
			packageVersionManifestPath := filepath.Join(packagePath, packageVersion, packageVersionManifestName)

			ok, err := svc.fs.Exists(packageVersionManifestPath)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}

			if !evaluate(where, packageName, packageVersion) {
				continue
			}

			version := &PackageVersion{
				Manifest: &PackageVersionManifest{
					Name: packageVersionManifestName,
				},
				Version: packageVersion,
				Feeds:   []string{svc.name}}

			if include != nil && include.Content {
				content, err := svc.fs.Read(packageVersionManifestPath)
				if err != nil {
					return nil, err
				}
				stringContent := string(content)
				version.Manifest.Content = &stringContent
			}
			versions = append(versions, version)
		}

		if len(versions) == 0 {
			continue
		}

		pkg := Package{
			Name:     packageName,
			Versions: versions,
		}
		packages = append(packages, &pkg)
	}
	return packages, nil
}

func (svc *fsFeedService) Create(request *FeedCreateRequest) (*FeedCreateResponse, error) {
	return nil, nil
}
