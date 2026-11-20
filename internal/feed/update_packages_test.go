package feed_test

import (
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
	feedfs "github.com/patrickhuber/wrangle/internal/feed/fs"
	"github.com/patrickhuber/wrangle/internal/resource"
	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
	"github.com/stretchr/testify/require"
)

func TestUpdatePackages(t *testing.T) {
	t.Run("generates new version from github release", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		logger := log.Memory()
		rootDir := "/opt/wrangle"
		workDir := path.Join(rootDir, "feed")
		// Set up an item with resource config, state and template
		itemName := "testpkg"
		itemDir := workDir + "/" + itemName

		require.NoError(t, fs.MkdirAll(itemDir, 0755))
		require.NoError(t, fs.MkdirAll(path.Join(itemDir, "0.6"), 0755))
		require.NoError(t, fs.WriteFile(itemDir+"/resource.yml", []byte(`resource:
  name: github-release
  type: github-release
  source:
    owner: testowner
    repository: testrepo
    version_regex: '[0-9]+\.[0-9]+\.[0-9]+'
    access_token: test-token
`), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/template.yml", []byte(`package:
  name: testpkg
  version: {{ .version }}
  targets:
  {{ range .platforms -}}
  - platform: {{ .platform }}
    architecture: {{ index .architectures 0 }}
    steps:
    - action: download
      with:
        url: https://example.com/testpkg-{{ $.version }}.tar.gz
  {{ end -}}
`), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/platforms.yml", []byte(`platforms:
- platform: linux
  architectures:
  - amd64
`), 0644))

		// Set up fake GitHub releases (v1.0.0 already exists, v1.1.0 is new)
		v100 := "v1.0.0"
		v110 := "v1.1.0"
		fakeReleases := []*github.RepositoryRelease{
			{TagName: &v110},
			{TagName: &v100},
		}
		fakeGH := &githubrelease.FakeGitHub{Releases: fakeReleases}

		factory := feedfs.NewLocalServiceFactory(fs, path, logger)

		// Use a fake GitHub client factory that returns our fake client
		ghProvider := githubrelease.NewProviderWithClientFactory(func(token string) githubrelease.GitHub {
			require.Equal(t, "test-token", token)
			return fakeGH
		}, logger)
		resourceFactory := resource.NewFactory(ghProvider)

		updateSvc := feed.NewUpdatePackages(factory, resourceFactory, logger)

		resp, err := updateSvc.Execute(&feed.UpdatePackagesRequest{Directory: rootDir})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 1, resp.Updated)

		// Verify the new version file was created
		versionFile := workDir + "/testpkg/1.1.0/package.yml"
		exists, err := fs.Exists(versionFile)
		require.NoError(t, err)
		require.True(t, exists, "version file should exist at %s", versionFile)
		exists, err = fs.Exists(workDir + "/testpkg/1.0.0/package.yml")
		require.NoError(t, err)
		require.False(t, exists, "only the latest version should be generated without state")

		// Verify state was updated
		stateData, err := fs.ReadFile(itemDir + "/state.yml")
		require.NoError(t, err)
		require.Contains(t, string(stateData), "1.1.0")
	})

	t.Run("uses the registered resource provider", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		logger := log.Memory()
		rootDir := "/opt/wrangle"
		workDir := path.Join(rootDir, "feed")
		itemDir := workDir + "/testpkg"

		require.NoError(t, fs.MkdirAll(itemDir, 0755))
		require.NoError(t, fs.WriteFile(itemDir+"/resource.yml", []byte(`resource:
  name: test-resource
  type: test-resource
  source: {}
`), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/state.yml", []byte("version: 1.0.0\n"), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/template.yml", []byte(`package:
  name: testpkg
  version: {{ .version }}
  targets:
  {{ range .platforms -}}
  - platform: {{ .platform }}
    architecture: {{ index .architectures 0 }}
  {{ end -}}
`), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/platforms.yml", []byte(`platforms:
- platform: linux
  architectures:
  - amd64
`), 0644))

		feedFactory := feedfs.NewLocalServiceFactory(fs, path, logger)
		resourceFactory := resource.NewFactory(fakeResourceProvider{})
		updateSvc := feed.NewUpdatePackages(feedFactory, resourceFactory, logger)

		resp, err := updateSvc.Execute(&feed.UpdatePackagesRequest{Directory: rootDir})
		require.NoError(t, err)
		require.Equal(t, 1, resp.Updated)

		exists, err := fs.Exists(itemDir + "/1.1.0/package.yml")
		require.NoError(t, err)
		require.True(t, exists)
	})

	t.Run("skips item without resource config", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		logger := log.Memory()
		rootDir := "/opt/wrangle"
		workDir := path.Join(rootDir, "feed")
		// Set up an item without a resource config
		itemName := "noresource"
		itemDir := workDir + "/" + itemName
		require.NoError(t, fs.MkdirAll(itemDir, 0755))
		require.NoError(t, fs.WriteFile(itemDir+"/state.yml", []byte("version: 1.0.0\n"), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/template.yml", []byte(""), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/platforms.yml", []byte(""), 0644))

		factory := feedfs.NewLocalServiceFactory(fs, path, logger)
		resourceFactory := resource.NewFactory()

		updateSvc := feed.NewUpdatePackages(factory, resourceFactory, logger)
		resp, err := updateSvc.Execute(&feed.UpdatePackagesRequest{Directory: rootDir})
		require.NoError(t, err)
		require.Equal(t, 0, resp.Updated)
	})
}

type fakeResourceProvider struct{}

func (fakeResourceProvider) Type() string {
	return "test-resource"
}

func (fakeResourceProvider) Create(resource.Config) (resource.Resource, error) {
	return fakeResource{}, nil
}

type fakeResource struct{}

func (fakeResource) Type() string {
	return "test-resource"
}

func (fakeResource) Name() string {
	return "test-resource"
}

func (fakeResource) Versions(string) ([]string, error) {
	return []string{"1.0.0", "1.1.0"}, nil
}
