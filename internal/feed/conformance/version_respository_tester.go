package conformance

import (
	"testing"

	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/patrickhuber/wrangle/internal/packages"
	"github.com/stretchr/testify/require"
)

type VersionRepositoryTester interface {
	CanGetSingleVersion(t *testing.T)
	CanListAllVersions(t *testing.T)
	CanAddVersion(t *testing.T, packageName, version string)
	CanUpdateVersionNumber(t *testing.T, packageName string, version string, newVersion string)
	CanAddTask(t *testing.T)
}

type versionRepositoryTester struct {
	repo feed.VersionRepository
}

func NewVersionRepositoryTester(repo feed.VersionRepository) VersionRepositoryTester {
	return &versionRepositoryTester{
		repo: repo,
	}
}

func (test *versionRepositoryTester) CanGetSingleVersion(t *testing.T) {
	packageName := "test"
	version := "1.0.0"
	v, err := test.repo.Get(packageName, version)
	require.NoError(t, err)
	require.NotNil(t, v)
}

func (test *versionRepositoryTester) CanListAllVersions(t *testing.T) {
	packageName := "test"
	expectedCount := 3
	v, err := test.repo.List(packageName)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, expectedCount, len(v))
}

func (test *versionRepositoryTester) CanAddVersion(t *testing.T, packageName, version string) {
	v := &packages.Version{
		Version: version,
		Manifest: &packages.Manifest{
			Package: &packages.ManifestPackage{
				Name:    packageName,
				Version: version,
				Targets: []*packages.ManifestTarget{
					{
						Platform:     platform.Linux,
						Architecture: arch.AMD64,
						Steps:        []*packages.ManifestStep{},
					},
				},
			},
		},
	}
	err := test.repo.Save(packageName, v)
	require.NoError(t, err)
	v, err = test.repo.Get(packageName, version)
	require.NoError(t, err)
	require.Equal(t, version, v.Version)
}

func (test *versionRepositoryTester) CanUpdateVersionNumber(t *testing.T, packageName string, version string, newVersion string) {

	v, err := test.repo.Get(packageName, version)
	require.NoError(t, err)
	require.NotNil(t, v)

	v.Version = newVersion
	v.Manifest.Package.Version = newVersion

	err = test.repo.Save(packageName, v)
	require.NoError(t, err)

	v, err = test.repo.Get(packageName, newVersion)
	require.NoError(t, err)
	require.Equal(t, newVersion, v.Version)
}

func (test *versionRepositoryTester) CanAddTask(t *testing.T) {
	packageName := "test"
	version := "1.0.0"

	v, err := test.repo.Get(packageName, version)
	require.NoError(t, err)

	require.NotNil(t, v)
	require.Equal(t, 3, len(v.Manifest.Package.Targets))

	task := &packages.ManifestStep{
		Action: "test",
		With:   map[string]any{},
	}
	target := v.Manifest.Package.Targets[0]
	target.Steps = append(v.Manifest.Package.Targets[0].Steps, task)

	err = test.repo.Save(packageName, v)
	require.NoError(t, err)

	v, err = test.repo.Get(packageName, version)
	require.NoError(t, err)
	require.Equal(t, 3, len(v.Manifest.Package.Targets))
	require.Equal(t, 2, len(v.Manifest.Package.Targets[0].Steps))
}
