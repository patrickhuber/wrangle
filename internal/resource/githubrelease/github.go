package githubrelease

import (
	"context"

	"github.com/google/go-github/v62/github"
)

// GitHub defines a wrapper to calling the github client
type GitHub interface {
	ListReleases(request *ListRequest) ([]*github.RepositoryRelease, error)
}

type githubClient struct {
	client *github.Client
}

// ListRequest defines the list parameters
type ListRequest struct {
	Owner      string
	Repository string
}

func (g *githubClient) ListReleases(request *ListRequest) ([]*github.RepositoryRelease, error) {
	listOptions := &github.ListOptions{PerPage: 100}

	// loop through all pages of releases and accumulate into the releases structure
	releases := []*github.RepositoryRelease{}
	for {
		// fetch the page
		releasesPage, response, err := g.client.Repositories.ListReleases(
			context.Background(),
			request.Owner,
			request.Repository,
			listOptions)

		if err != nil {
			return nil, err
		}

		// append to all releases
		releases = append(releases, releasesPage...)

		// set the next page
		if response.NextPage > 0 {
			listOptions.Page = response.NextPage
			continue
		}

		// no more pages
		err = response.Body.Close()
		if err != nil {
			return nil, err
		}
		break
	}

	return releases, nil
}
