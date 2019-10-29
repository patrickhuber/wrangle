package packages

import (
	"fmt"
	"strings"

	semver "github.com/hashicorp/go-version"
	"github.com/patrickhuber/wrangle/filepath"
	"github.com/patrickhuber/wrangle/filesystem"
	"github.com/patrickhuber/wrangle/settings"
)

type fsContextProvider struct {
	fs    filesystem.FileSystem
	paths *settings.Paths
}

// NewFsContextProvider creates a context provider for the file system
func NewFsContextProvider(fs filesystem.FileSystem, paths *settings.Paths) ContextProvider {
	return &fsContextProvider{
		fs:    fs,
		paths: paths,
	}
}

func (p *fsContextProvider) Get(packageName, packageVersion string) (PackageContext, error) {
	// packagesDirectory
	//   ex: /packages
	// packagePath
	//   ex: /packages/test
	// packageVersionPath
	//   ex: /packages/test/1.0.0/
	// packageVersionManifestPath
	//   ex: /packages/test/1.0.0/test.1.0.0.yml
	packagePath, err := p.getPackagePath(p.paths.Packages, packageName)
	if err != nil {
		return nil, err
	}

	packageVersion, err = p.getPackageVersion(packagePath, packageVersion)
	if err != nil {
		return nil, err
	}

	packageVersionPath := fmt.Sprintf("%s/%s", packagePath, packageVersion)
	packageVersionManifestPath := fmt.Sprintf("%s/%s.%s.yml", packageVersionPath, packageName, packageVersion)

	packageContext := NewContext(
		p.paths,
		packagePath,
		packageVersionPath,
		packageVersionManifestPath)
	return packageContext, nil
}

func (p *fsContextProvider) getPackagePath(packagesRoot, packageName string) (string, error) {
	if strings.TrimSpace(packageName) == "" {
		return "", fmt.Errorf("package name is required")
	}

	packagePath := filepath.Join(packagesRoot, packageName)
	return packagePath, nil
}

func (p *fsContextProvider) getPackageVersion(packagePath, packageVersion string) (string, error) {
	useLatestVersion := len(strings.TrimSpace(packageVersion)) == 0

	if !useLatestVersion {
		return packageVersion, nil
	}

	return p.findLatestPackageVersion(packagePath)
}

func (p *fsContextProvider) findLatestPackageVersion(packagePath string) (string, error) {
	files, err := p.fs.ReadDir(packagePath)
	if err != nil {
		return "", err
	}

	var latest *semver.Version

	for _, file := range files {

		if !file.IsDir() {
			continue
		}

		version := file.Name()
		v, err := semver.NewVersion(version)
		if err != nil {
			continue
		}

		if latest == nil {
			latest = v
			continue
		}

		if v.GreaterThan(latest) {
			latest = v
		}
	}

	if latest == nil {
		return "", fmt.Errorf("unable to determine latest version of package")
	}

	return latest.String(), nil
}
