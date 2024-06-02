package githubrelease

import "github.com/google/go-github/v62/github"

type FakeGitHub struct {
	Releases []*github.RepositoryRelease
}

func (g *FakeGitHub) ListReleases(request *ListRequest) ([]*github.RepositoryRelease, error) {
	return g.Releases, nil
}
