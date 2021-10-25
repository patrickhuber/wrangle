package githubrelease

import "github.com/patrickhuber/wrangle/pkg/resource"

type githubReleaseResource struct {
	name string
	svc  Service
}

func NewResource(name string, svc Service) resource.Resource {
	return &githubReleaseResource{
		name: name,
		svc:  svc,
	}
}

func (r *githubReleaseResource) Type() string {
	return "githubrelease"
}

func (r *githubReleaseResource) Name() string {
	return r.name
}

func (r *githubReleaseResource) Versions() ([]string, error) {
	request := &CheckRequest{}
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
