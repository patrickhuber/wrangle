package githubrelease_test

import (
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/resource"
	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
	"github.com/stretchr/testify/require"
)

func TestProviderValidatesSource(t *testing.T) {
	tests := []struct {
		name   string
		source map[string]string
		err    string
	}{
		{
			name: "missing owner",
			source: map[string]string{
				"repository": "testrepo",
			},
			err: `"owner"`,
		},
		{
			name: "missing repository",
			source: map[string]string{
				"owner": "testowner",
			},
			err: `"repository"`,
		},
		{
			name: "invalid version regex",
			source: map[string]string{
				"owner":         "testowner",
				"repository":    "testrepo",
				"version_regex": "[",
			},
			err: `"version_regex"`,
		},
		{
			name: "invalid tag filter",
			source: map[string]string{
				"owner":      "testowner",
				"repository": "testrepo",
				"tag_filter": "[",
			},
			err: `"tag_filter"`,
		},
		{
			name: "unknown key",
			source: map[string]string{
				"owner":      "testowner",
				"repository": "testrepo",
				"unknown":    ".*",
			},
			err: `unknown github-release source key "unknown"`,
		},
	}

	provider := githubrelease.NewProviderWithClientFactory(func(string) githubrelease.GitHub {
		return &githubrelease.FakeGitHub{}
	}, log.Memory())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := provider.Create(resource.Config{
				Name:   "test",
				Type:   "github-release",
				Source: tt.source,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.err)
		})
	}
}

func TestProviderAcceptsLegacyVersionRegex(t *testing.T) {
	first := "1.0.0"
	second := "2.0.0"
	provider := githubrelease.NewProviderWithClientFactory(func(string) githubrelease.GitHub {
		return &githubrelease.FakeGitHub{
			Releases: []*github.RepositoryRelease{
				{TagName: &first},
				{TagName: &second},
			},
		}
	}, log.Memory())

	resource, err := provider.Create(resource.Config{
		Name: "test",
		Type: "github-release",
		Source: map[string]string{
			"owner":         "testowner",
			"repository":    "testrepo",
			"version-regex": `2\..*`,
		},
	})
	require.NoError(t, err)

	versions, err := resource.Versions("")
	require.NoError(t, err)
	require.Equal(t, []string{"2.0.0"}, versions)
}
