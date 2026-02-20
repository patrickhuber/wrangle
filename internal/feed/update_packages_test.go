package feed_test

import (
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/feed"
	feedfs "github.com/patrickhuber/wrangle/internal/feed/fs"
	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
	"github.com/stretchr/testify/require"
)

func TestUpdatePackages(t *testing.T) {
	t.Run("generates new version from github release", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		logger := log.Memory()
		workDir := "/opt/wrangle/feed"
		svc := feedfs.NewService("test", fs, path, workDir, logger)

		// Set up an item with resource config, state and template
		itemName := "testpkg"
		itemDir := workDir + "/" + itemName

		require.NoError(t, fs.MkdirAll(itemDir, 0755))
		require.NoError(t, fs.WriteFile(itemDir+"/resource.yml", []byte(`resource:
  name: github-release
  type: github-release
  source:
    owner: testowner
    repository: testrepo
    version-regex: '[0-9]+\.[0-9]+\.[0-9]+'
`), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/state.yml", []byte("version: 1.0.0\n"), 0644))
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
		require.NoError(t, fs.WriteFile(itemDir+"/platforms.yml", []byte(`- platform: linux
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

		// Create a fake service factory that returns our fs-backed service
		cfg := config.Config{
			Spec: config.Spec{
				Feeds: []config.Feed{
					{Name: "test", Type: "test"},
				},
			},
		}
		cfgSvc := &fakeConfigService{cfg: cfg}
		factory := &fakeFeedServiceFactory{svc: svc}

		// Use a fake GitHub client factory that returns our fake client
		ghFactory := feed.GitHubClientFactory(func(token string) githubrelease.GitHub {
			return fakeGH
		})

		updateSvc := feed.NewUpdatePackagesWithGitHubClientFactory(factory, cfgSvc, ghFactory)

		resp, err := updateSvc.Execute(&feed.UpdatePackagesRequest{})
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 1, resp.Updated)

		// Verify the new version file was created
		versionFile := workDir + "/testpkg/1.1.0/package.yml"
		exists, err := fs.Exists(versionFile)
		require.NoError(t, err)
		require.True(t, exists, "version file should exist at %s", versionFile)

		// Verify state was updated
		stateData, err := fs.ReadFile(itemDir + "/state.yml")
		require.NoError(t, err)
		require.Contains(t, string(stateData), "1.1.0")
	})

	t.Run("skips item without resource config", func(t *testing.T) {
		h := cross.NewTest(platform.Linux, arch.AMD64)
		fs := h.FS()
		path := h.Path()
		logger := log.Memory()
		workDir := "/opt/wrangle/feed"
		svc := feedfs.NewService("test", fs, path, workDir, logger)

		// Set up an item without a resource config
		itemName := "noresource"
		itemDir := workDir + "/" + itemName
		require.NoError(t, fs.MkdirAll(itemDir, 0755))
		require.NoError(t, fs.WriteFile(itemDir+"/state.yml", []byte("version: 1.0.0\n"), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/template.yml", []byte(""), 0644))
		require.NoError(t, fs.WriteFile(itemDir+"/platforms.yml", []byte(""), 0644))

		cfg := config.Config{
			Spec: config.Spec{
				Feeds: []config.Feed{{Name: "test", Type: "test"}},
			},
		}
		cfgSvc := &fakeConfigService{cfg: cfg}
		factory := &fakeFeedServiceFactory{svc: svc}
		ghFactory := feed.GitHubClientFactory(func(token string) githubrelease.GitHub {
			return &githubrelease.FakeGitHub{}
		})

		updateSvc := feed.NewUpdatePackagesWithGitHubClientFactory(factory, cfgSvc, ghFactory)
		resp, err := updateSvc.Execute(&feed.UpdatePackagesRequest{})
		require.NoError(t, err)
		require.Equal(t, 0, resp.Updated)
	})
}

// fakeConfigService implements config.Service for testing
type fakeConfigService struct {
	cfg config.Config
}

func (f *fakeConfigService) Get() (config.Config, error) {
	return f.cfg, nil
}

// fakeFeedServiceFactory returns a predetermined feed.Service for testing
type fakeFeedServiceFactory struct {
	svc feed.Service
}

func (f *fakeFeedServiceFactory) Create(cfg config.Feed) (feed.Service, error) {
	return f.svc, nil
}
