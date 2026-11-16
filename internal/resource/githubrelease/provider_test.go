package githubrelease_test

import (
	"testing"

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
				"owner":         "testowner",
				"repository":    "testrepo",
				"version-regex": ".*",
			},
			err: `unknown github-release source key "version-regex"`,
		},
	}

	provider := githubrelease.NewProviderWithClientFactory(func(string) githubrelease.GitHub {
		return &githubrelease.FakeGitHub{}
	})

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
