package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/filepath"
	"github.com/spf13/afero"
)

type fsFeedService struct {
	fs   afero.Fs
	path string
}

// NewFsFeedService defines a feed service over the filesystem
func NewFsFeedService(fs afero.Fs, path string) FeedService {
	return &fsFeedService{
		fs:   fs,
		path: path,
	}
}

func (s *fsFeedService) List(request *FeedListRequest) (*FeedListResponse, error) {
	packageFolders, err := afero.ReadDir(s.fs, s.path)
	if err != nil {
		return nil, err
	}
	packages := []FeedListResponsePackage{}
	for _, packageFolder := range packageFolders {
		if !packageFolder.IsDir() {
			continue
		}
		packageName := packageFolder.Name()
		packagePath := filepath.Join(s.path, packageFolder.Name())
		packageVersions, err := afero.ReadDir(s.fs, packagePath)
		if err != nil {
			return nil, err
		}

		versions := []string{}
		for _, packageVersionFolder := range packageVersions {
			if !packageVersionFolder.IsDir() {
				continue
			}
			packageVersion := packageVersionFolder.Name()
			packageVersionManifestName := fmt.Sprintf("%s.%s.yml", packageName, packageVersion)
			packageVersionManifestPath := filepath.Join(packagePath, packageVersion, packageVersionManifestName)

			ok, err := afero.Exists(s.fs, packageVersionManifestPath)
			if err != nil {
				return nil, err
			}
			if !ok {
				continue
			}
			versions = append(versions, packageVersion)
		}

		if len(versions) == 0 {
			continue
		}

		pkg := FeedListResponsePackage{
			Name:     packageName,
			Versions: versions,
		}
		packages = append(packages, pkg)
	}
	return &FeedListResponse{
		Packages: packages,
	}, nil
}
