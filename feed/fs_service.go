package feed

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
)

type fsService struct {
	fs   filesystem.FileSystem
	path string
	name string
}

// NewFsService defines a feed service over the filesystem
func NewFsService(fs filesystem.FileSystem, path string) Service {
	return &fsService{
		fs:   fs,
		path: path,
		name: "local",
	}
}

func (svc *fsService) List(request *ListRequest) (*ListResponse, error) {
	packages, err := svc.find(nil, nil)
	if err != nil {
		return nil, err
	}
	return &ListResponse{
		Packages: packages,
	}, nil
}

func (svc *fsService) Get(request *GetRequest) (*GetResponse, error) {
	versions := []string{}
	if strings.TrimSpace(request.Version) != "" {
		versions = append(versions, request.Version)
	}
	where := &packageCriteriaWhere{
		Or: []*packageCriteriaAnd{
			{
				And: []*packageCriteria{
					{
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

	return &GetResponse{
		Package: pkg,
	}, nil
}

func (svc *fsService) Info(request *InfoRequest) (*InfoResponse, error) {
	return &InfoResponse{
		URI: svc.path,
	}, nil
}

func (svc *fsService) find(where *packageCriteriaWhere, include *packageInclude) ([]*Package, error) {
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

		latest, err := svc.getLatestTag(packagePath)
		if err != nil {
			return nil, err
		}

		ignoreLatest := latest != nil

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
				version.Manifest.Content = stringContent
			}

			ver, err := semver.NewVersion(version.Version)
			if err != nil {
				return nil, err
			}

			if latest == nil || (latest.Compare(ver) == -1 && !ignoreLatest) {
				latest = ver
			}

			versions = append(versions, version)
		}

		if len(versions) == 0 {
			continue
		}

		pkg := Package{
			Name:     packageName,
			Versions: versions,
			Latest:   latest.String(),
		}
		packages = append(packages, &pkg)
	}
	return packages, nil
}

func (svc *fsService) getLatestTag(packagePath string) (*semver.Version, error) {
	latestTagPath := filepath.Join(packagePath, "latest")
	tagExists, err := svc.fs.Exists(latestTagPath)

	if err != nil {
		return nil, err
	}

	if !tagExists {
		return nil, nil
	}

	version, err := svc.fs.Read(latestTagPath)
	if err != nil {
		return nil, err
	}

	return semver.NewVersion(string(version))
}

func (svc *fsService) Create(request *CreateRequest) (*CreateResponse, error) {
	return nil, fmt.Errorf("not implemented")
}

func (svc *fsService) Latest(request *LatestRequest) (*LatestResponse, error) {
	response, err := svc.Get(&GetRequest{Name: request.Name})
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
	return &LatestResponse{
		Package: response.Package,
	}, nil
}
