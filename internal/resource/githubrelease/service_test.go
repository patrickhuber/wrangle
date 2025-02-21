package githubrelease_test

import (
	"testing"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	setup := func(versions []string) githubrelease.Service {
		releases := []*github.RepositoryRelease{}
		for i := range versions {
			version := versions[i]
			releases = append(releases, &github.RepositoryRelease{TagName: &version})
		}
		client := &githubrelease.FakeGitHub{Releases: releases}
		return githubrelease.NewService(client)
	}
	t.Run("check when releases", func(t *testing.T) {
		versions := []string{"v1.0.0", "v2.0.0"}
		service := setup(versions)
		request := &githubrelease.CheckRequest{}
		response, err := service.Check(request)
		require.NoError(t, err)
		require.NotNil(t, response.Versions)
		require.Equal(t, 2, len(response.Versions))
	})
	t.Run("returns greater", func(t *testing.T) {
		versions := []string{"v1.10.0", "v1.10.1", "v1.11.0", "v1.11.1"}

		service := setup(versions)
		request := &githubrelease.CheckRequest{
			Version: githubrelease.Version{
				ID: "1.10.1",
			},
		}
		response, err := service.Check(request)
		require.NoError(t, err)
		require.NotNil(t, response.Versions)
		require.Equal(t, 3, len(response.Versions))
	})
	t.Run("returns current release", func(t *testing.T) {
		versions := []string{"v1.10.0", "v1.10.1", "v1.11.0", "v1.11.1"}
		service := setup(versions)
		request := &githubrelease.CheckRequest{
			Version: githubrelease.Version{
				ID: "1.11.1",
			},
		}
		response, err := service.Check(request)
		require.NoError(t, err)
		require.NotNil(t, response.Versions)
		require.Equal(t, 1, len(response.Versions))
		require.Equal(t, "1.11.1", response.Versions[0].ID)
	})
}
