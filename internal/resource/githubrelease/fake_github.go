package githubrelease

import (
	"github.com/google/go-github/v62/github"
)

type FakeGitHub struct {
	Releases []*github.RepositoryRelease
	Err      error
}

func (g *FakeGitHub) ListReleases(request *ListRequest) ([]*github.RepositoryRelease, error) {
	return g.Releases, g.Err
}
