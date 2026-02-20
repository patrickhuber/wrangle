package feed

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/patrickhuber/wrangle/internal/config"
	"github.com/patrickhuber/wrangle/internal/resource/githubrelease"
)

// UpdatePackagesRequest contains options for the feed update command
type UpdatePackagesRequest struct{}

// UpdatePackagesResponse contains the result of updating feed packages
type UpdatePackagesResponse struct {
	Updated int
}

// UpdatePackages updates feed packages to their latest versions based on resource configurations
type UpdatePackages interface {
	Execute(r *UpdatePackagesRequest) (*UpdatePackagesResponse, error)
}

// GitHubClientFactory creates a GitHub client, optionally authenticated with a token
type GitHubClientFactory func(token string) githubrelease.GitHub

type updatePackages struct {
	serviceFactory      ServiceFactory
	configuration       config.Service
	gitHubClientFactory GitHubClientFactory
}

// NewUpdatePackages creates a new UpdatePackages service using the default GitHub client
func NewUpdatePackages(serviceFactory ServiceFactory, configuration config.Service) UpdatePackages {
	return &updatePackages{
		serviceFactory:      serviceFactory,
		configuration:       configuration,
		gitHubClientFactory: githubrelease.NewGitHubClient,
	}
}

// NewUpdatePackagesWithGitHubClientFactory creates an UpdatePackages with a custom GitHub client factory (useful for testing)
func NewUpdatePackagesWithGitHubClientFactory(serviceFactory ServiceFactory, configuration config.Service, factory GitHubClientFactory) UpdatePackages {
	return &updatePackages{
		serviceFactory:      serviceFactory,
		configuration:       configuration,
		gitHubClientFactory: factory,
	}
}

func (u *updatePackages) Execute(r *UpdatePackagesRequest) (*UpdatePackagesResponse, error) {
	cfg, err := u.configuration.Get()
	if err != nil {
		return nil, err
	}
	if len(cfg.Spec.Feeds) == 0 {
		return &UpdatePackagesResponse{}, nil
	}

	totalUpdated := 0
	for _, f := range cfg.Spec.Feeds {
		feedSvc, err := u.serviceFactory.Create(f)
		if err != nil {
			return nil, err
		}

		updated, err := u.updateFeed(feedSvc)
		if err != nil {
			return nil, err
		}
		totalUpdated += updated
	}

	return &UpdatePackagesResponse{Updated: totalUpdated}, nil
}

func (u *updatePackages) updateFeed(feedSvc Service) (int, error) {
	items, err := feedSvc.List(&ListRequest{})
	if err != nil {
		return 0, err
	}

	updated := 0
	for _, item := range items.Items {
		n, err := u.updateItem(feedSvc, item)
		if err != nil {
			return updated, err
		}
		updated += n
	}
	return updated, nil
}

func (u *updatePackages) updateItem(feedSvc Service, item *Item) (int, error) {
	if item.Resource == nil {
		return 0, nil
	}
	if item.Resource.Type != "github-release" {
		return 0, nil
	}

	source := buildSource(item.Resource)
	currentVersion := ""
	if item.State != nil {
		currentVersion = item.State.LatestVersion
	}

	client := u.gitHubClientFactory(source.AccessToken)
	svc := githubrelease.NewService(client)

	checkResp, err := svc.Check(&githubrelease.CheckRequest{
		Source:  source,
		Version: githubrelease.Version{ID: currentVersion},
	})
	if err != nil {
		return 0, fmt.Errorf("error checking versions for %s: %w", item.Package.Name, err)
	}

	// collect versions newer than current
	newVersions := []string{}
	latestVersion := currentVersion
	for _, v := range checkResp.Versions {
		if v.ID == currentVersion {
			continue
		}
		newVersions = append(newVersions, v.ID)
		latestVersion = v.ID
	}

	if len(newVersions) == 0 {
		return 0, nil
	}

	// render and save each new version
	platformData := toPlatformData(item.Platforms)
	for _, version := range newVersions {
		rendered, err := renderTemplate(item.Template, version, platformData)
		if err != nil {
			return 0, fmt.Errorf("error rendering template for %s@%s: %w", item.Package.Name, version, err)
		}
		if err := feedSvc.SaveVersion(item.Package.Name, version, rendered); err != nil {
			return 0, fmt.Errorf("error saving version %s@%s: %w", item.Package.Name, version, err)
		}
	}

	// update state with latest version
	if latestVersion != currentVersion {
		_, err = feedSvc.Update(&UpdateRequest{
			Items: &ItemUpdate{
				Modify: []*ItemModify{
					{
						Name: item.Package.Name,
						State: &StateModify{
							LatestVersion: &latestVersion,
						},
					},
				},
			},
		})
		if err != nil {
			return 0, fmt.Errorf("error updating state for %s: %w", item.Package.Name, err)
		}
	}

	return len(newVersions), nil
}

func buildSource(resource *ResourceConfig) githubrelease.Source {
	return githubrelease.Source{
		Owner:        resource.Source["owner"],
		Repository:   resource.Source["repository"],
		VersionRegex: resource.Source["version-regex"],
		TagFilter:    resource.Source["tag_filter"],
		AccessToken:  resource.Source["access_token"],
	}
}

// toPlatformData converts feed.Platform structs to template-compatible map data
func toPlatformData(platforms []*Platform) []map[string]any {
	result := make([]map[string]any, 0, len(platforms))
	for _, p := range platforms {
		result = append(result, map[string]any{
			"platform":      p.Name,
			"architectures": p.Architectures,
		})
	}
	return result
}

// renderTemplate executes the item template with the given version and platform data
func renderTemplate(templateStr string, version string, platforms []map[string]any) ([]byte, error) {
	tmpl, err := template.New("package").Parse(templateStr)
	if err != nil {
		return nil, err
	}
	data := map[string]any{
		"version":   version,
		"platforms": platforms,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
