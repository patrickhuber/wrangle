package git

import (
	"io"

	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/patrickhuber/wrangle/pkg/feed"
	"github.com/patrickhuber/wrangle/pkg/packages"
)

type packageVersionRepository struct {
	tree *object.Tree
}

type PackageVersionManifest struct {
	Package *packages.PackageVersion `yaml: "package", json: "package"`
}

func (r *packageVersionRepository) Get(packageName string, version string) (*packages.PackageVersion, error) {
	packageTree, err := r.tree.Tree(packageName)
	if err != nil {
		return nil, err
	}

	return r.get(packageTree, version)
}

func (r *packageVersionRepository) get(packageTree *object.Tree, version string) (*packages.PackageVersion, error) {
	packageVersionTree, err := packageTree.Tree(version)
	if err != nil {
		return nil, err
	}

	manifest := &PackageVersionManifest{}
	err = DecodeYamlFileFromGitTree(packageVersionTree, "package.yml", manifest)
	if err != nil {
		return nil, err
	}

	return &packages.PackageVersion{
		Version: manifest.Package.Version,
		Targets: manifest.Package.Targets,
	}, nil
}

func (r *packageVersionRepository) List(packageName string, latestVersion string, expand *feed.ItemReadExpandPackage) ([]*packages.PackageVersion, error) {

	packageTree, err := r.tree.Tree(packageName)
	if err != nil {
		return nil, err
	}

	seen := map[plumbing.Hash]bool{}
	walker := object.NewTreeWalker(packageTree, false, seen)
	versions := []*packages.PackageVersion{}

	for {
		name, entry, err := walker.Next()
		if err == io.EOF {
			break
		}
		if entry.Mode.IsFile() {
			continue
		}

		isMatch := expand.IsMatch(name, latestVersion)
		if !isMatch {
			continue
		}
		version, err := r.get(packageTree, name)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}
