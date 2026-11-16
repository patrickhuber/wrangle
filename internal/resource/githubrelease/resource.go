package githubrelease

import "github.com/patrickhuber/wrangle/internal/resource"

type githubReleaseResource struct {
	name   string
	source Source
	svc    Service
}

func NewResource(name string, source Source, svc Service) resource.Resource {
	return &githubReleaseResource{
		name:   name,
		source: source,
		svc:    svc,
	}
}

func (r *githubReleaseResource) Type() string {
	return "github-release"
}

func (r *githubReleaseResource) Name() string {
	return r.name
}

func (r *githubReleaseResource) Versions(currentVersion string) ([]string, error) {
	request := &CheckRequest{
		Source:  r.source,
		Version: Version{ID: currentVersion},
	}
	response, err := r.svc.Check(request)
	if err != nil {
		return nil, err
	}
	versions := []string{}
	for _, v := range response.Versions {
		versions = append(versions, v.ID)
	}
	return versions, nil
}
