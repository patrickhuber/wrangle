package feed

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/semver/v3"
	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/resource"
)

// UpdatePackagesRequest contains options for the feed update command
type UpdatePackagesRequest struct {
	Directory string
}

// UpdatePackagesResponse contains the result of updating feed packages
type UpdatePackagesResponse struct {
	Updated int
}

// UpdatePackages updates feed packages to their latest versions based on resource configurations
type UpdatePackages interface {
	Execute(r *UpdatePackagesRequest) (*UpdatePackagesResponse, error)
}

type updatePackages struct {
	serviceFactory  LocalServiceFactory
	resourceFactory resource.Factory
	logger          log.Logger
}

func NewUpdatePackages(serviceFactory LocalServiceFactory, resourceFactory resource.Factory, logger log.Logger) UpdatePackages {
	return &updatePackages{
		serviceFactory:  serviceFactory,
		resourceFactory: resourceFactory,
		logger:          logger,
	}
}

func (u *updatePackages) Execute(r *UpdatePackagesRequest) (*UpdatePackagesResponse, error) {
	feedSvc, err := u.serviceFactory.Create(r.Directory)
	if err != nil {
		return nil, err
	}

	updated, err := u.updateFeed(feedSvc)
	if err != nil {
		return nil, err
	}

	return &UpdatePackagesResponse{Updated: updated}, nil
}

func (u *updatePackages) updateFeed(feedSvc Service) (int, error) {
	items, err := feedSvc.List(&ListRequest{
		AllowMissingState:        true,
		SkipPackageVersions:      true,
		SkipItemsWithoutResource: true,
	})
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
		u.logger.Infof("skipping feed package %s: missing resource configuration", item.Package.Name)
		return 0, nil
	}

	currentVersion := ""
	hasState := item.State != nil
	if item.State != nil {
		currentVersion = item.State.LatestVersion
	}

	resource, err := u.resourceFactory.Create(resource.Config{
		Name:   item.Resource.Name,
		Type:   item.Resource.Type,
		Source: item.Resource.Source,
	})
	if err != nil {
		return 0, fmt.Errorf("create resource for %s: %w", item.Package.Name, err)
	}

	versions, err := resource.Versions(currentVersion)
	if err != nil {
		return 0, fmt.Errorf("error checking versions for %s: %w", item.Package.Name, err)
	}
	if len(versions) == 0 {
		u.logger.Infof("skipping feed package %s: no matching releases", item.Package.Name)
		return 0, nil
	}

	// collect versions newer than current, using semver to find the true latest
	newVersions := []string{}
	var latestSemver *semver.Version
	latestVersion := currentVersion
	for _, version := range versions {
		if version == currentVersion {
			continue
		}
		newVersions = append(newVersions, version)
		sv, err := semver.NewVersion(version)
		if err == nil && (latestSemver == nil || sv.GreaterThan(latestSemver)) {
			latestSemver = sv
			latestVersion = version
		}
	}

	if len(newVersions) == 0 {
		return 0, nil
	}
	if !hasState {
		newVersions = []string{latestVersion}
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
