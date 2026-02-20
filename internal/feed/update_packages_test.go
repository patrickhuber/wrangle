package feed_test

import (
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/feed"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

type fakeGitHubReleaseClient struct {
	latest   *github.RepositoryRelease
	pages    map[int][]*github.RepositoryRelease
	hasNext  map[int]bool
}

func (f *fakeGitHubReleaseClient) GetLatestRelease(owner, repository string) (*github.RepositoryRelease, error) {
	return f.latest, nil
}

func (f *fakeGitHubReleaseClient) ListReleases(owner, repository string, page, perPage int) ([]*github.RepositoryRelease, bool, error) {
	return f.pages[page], f.hasNext[page], nil
}

func TestUpdatePackagesGeneratesVersionsAndUpdatesState(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	logger := log.Memory()

	wd, err := target.OS().WorkingDirectory()
	require.NoError(t, err)

	feedDir := target.Path().Join(wd, "feed", "shim")
	err = target.FS().MkdirAll(feedDir, 0775)
	require.NoError(t, err)

	resource := `resource:
  type: github-release
  source:
    owner: patrickhuber
    repository: shim
    version-regex: '[0-9]+\.[0-9]+\.[0-9]+'
`
	platforms := `platforms:
- platform: linux
  architectures: [amd64]
`
	template := `package:
  name: shim
  version: {{ .version }}
  targets:
  {{- range .platforms }}
  - platform: {{ .platform }}
    architecture: {{ index .architectures 0 }}
    steps:
    - action: move
      with:
        source: shim
        destination: shim
  {{- end }}
`
	state := `version: 1.0.0
`

	err = target.FS().WriteFile(target.Path().Join(feedDir, "resource.yml"), []byte(resource), 0644)
	require.NoError(t, err)
	err = target.FS().WriteFile(target.Path().Join(feedDir, "platforms.yml"), []byte(platforms), 0644)
	require.NoError(t, err)
	err = target.FS().WriteFile(target.Path().Join(feedDir, "template.yml"), []byte(template), 0644)
	require.NoError(t, err)
	err = target.FS().WriteFile(target.Path().Join(feedDir, "state.yml"), []byte(state), 0644)
	require.NoError(t, err)

	version101 := "v1.0.1"
	version102 := "v1.0.2"
	client := &fakeGitHubReleaseClient{
		latest: &github.RepositoryRelease{TagName: &version102},
		pages: map[int][]*github.RepositoryRelease{
			1: {
				{TagName: &version102},
				{TagName: &version101},
				{TagName: stringPtr("v1.0.0")},
			},
		},
		hasNext: map[int]bool{1: false},
	}

	svc := feed.NewUpdatePackagesWithClient(target.FS(), target.Path(), target.OS(), logger, client)
	response, err := svc.Execute(&feed.UpdatePackagesRequest{})
	require.NoError(t, err)
	require.Equal(t, 1, response.Packages)
	require.Equal(t, 2, response.Versions)

	for _, version := range []string{"1.0.1", "1.0.2"} {
		manifestPath := target.Path().Join(feedDir, version, "package.yml")
		exists, err := target.FS().Exists(manifestPath)
		require.NoError(t, err)
		require.True(t, exists)
	}

	statePath := target.Path().Join(feedDir, "state.yml")
	stateContent, err := target.FS().ReadFile(statePath)
	require.NoError(t, err)

	actualState := &feed.State{}
	err = yaml.Unmarshal(stateContent, actualState)
	require.NoError(t, err)
	require.Equal(t, "1.0.2", actualState.LatestVersion)
}

func stringPtr(value string) *string {
	return &value
}
