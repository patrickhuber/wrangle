package feed

import (
	"bytes"
	"context"
	"regexp"
	"strings"
	"text/template"

	"github.com/google/go-github/v62/github"
	"github.com/patrickhuber/go-cross/filepath"
	"github.com/patrickhuber/go-cross/fs"
	"github.com/patrickhuber/go-cross/os"
	"github.com/patrickhuber/go-log"
	"gopkg.in/yaml.v3"
)

type UpdatePackages interface {
	Execute(request *UpdatePackagesRequest) (*UpdatePackagesResponse, error)
}

type UpdatePackagesRequest struct {
	FeedDirectory string
}

type UpdatePackagesResponse struct {
	Packages int
	Versions int
}

type GitHubReleaseClient interface {
	GetLatestRelease(owner, repository string) (*github.RepositoryRelease, error)
	ListReleases(owner, repository string, page, perPage int) ([]*github.RepositoryRelease, bool, error)
}

type updatePackages struct {
	fs      fs.FS
	path    filepath.Provider
	opsys   os.OS
	logger  log.Logger
	github  GitHubReleaseClient
	perPage int
}

func NewUpdatePackages(fs fs.FS, path filepath.Provider, opsys os.OS, logger log.Logger) UpdatePackages {
	return NewUpdatePackagesWithClient(fs, path, opsys, logger, &defaultGitHubReleaseClient{
		client: github.NewClient(nil),
	})
}

func NewUpdatePackagesWithClient(fs fs.FS, path filepath.Provider, opsys os.OS, logger log.Logger, github GitHubReleaseClient) UpdatePackages {
	return &updatePackages{
		fs:      fs,
		path:    path,
		opsys:   opsys,
		logger:  logger,
		github:  github,
		perPage: 100,
	}
}

type updateResourceFile struct {
	Resource updateResource `yaml:"resource"`
}

type updateResource struct {
	Type   string               `yaml:"type"`
	Source updateResourceSource `yaml:"source"`
}

type updateResourceSource struct {
	Owner        string `yaml:"owner"`
	Repository   string `yaml:"repository"`
	TagFilter    string `yaml:"tag_filter"`
	VersionRegex string `yaml:"version-regex"`
}

func (u *updatePackages) Execute(request *UpdatePackagesRequest) (*UpdatePackagesResponse, error) {
	if request == nil {
		request = &UpdatePackagesRequest{}
	}

	workingDirectory, err := u.opsys.WorkingDirectory()
	if err != nil {
		return nil, err
	}

	feedDirectory := strings.TrimSpace(request.FeedDirectory)
	if feedDirectory == "" {
		feedDirectory = "feed"
	}

	feedPath := u.path.Join(workingDirectory, feedDirectory)
	entries, err := u.fs.ReadDir(feedPath)
	if err != nil {
		return nil, err
	}

	response := &UpdatePackagesResponse{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		packageName := entry.Name()
		packageDir := u.path.Join(feedPath, packageName)
		updated, generated, err := u.updatePackage(packageName, packageDir)
		if err != nil {
			return nil, err
		}
		if updated {
			response.Packages++
		}
		response.Versions += generated
	}

	return response, nil
}

func (u *updatePackages) updatePackage(packageName, packageDir string) (bool, int, error) {
	resourcePath := u.path.Join(packageDir, "resource.yml")
	ok, err := u.fs.Exists(resourcePath)
	if err != nil {
		return false, 0, err
	}
	if !ok {
		u.logger.Infof("skip: %s, missing resource file '%s'", packageName, resourcePath)
		return false, 0, nil
	}

	resource := &updateResourceFile{}
	err = u.readYAML(resourcePath, resource)
	if err != nil {
		return false, 0, err
	}

	if resource.Resource.Type != "github-release" {
		u.logger.Infof("skip: %s, unknown release type '%s'", packageName, resource.Resource.Type)
		return false, 0, nil
	}

	owner := strings.TrimSpace(resource.Resource.Source.Owner)
	if owner == "" {
		u.logger.Infof("skip: %s, missing github-release.source.owner", packageName)
		return false, 0, nil
	}

	repository := strings.TrimSpace(resource.Resource.Source.Repository)
	if repository == "" {
		u.logger.Infof("skip: %s, missing github-release.source.repository", packageName)
		return false, 0, nil
	}

	versionRegex := strings.TrimSpace(resource.Resource.Source.VersionRegex)
	if versionRegex == "" {
		versionRegex = ".*"
	}

	statePath := u.path.Join(packageDir, "state.yml")
	stateExists, err := u.fs.Exists(statePath)
	if err != nil {
		return false, 0, err
	}

	latestRelease, err := u.github.GetLatestRelease(owner, repository)
	if err != nil {
		return false, 0, err
	}
	if latestRelease == nil || latestRelease.TagName == nil {
		u.logger.Infof("skip: %s, github latest release is missing a tag", packageName)
		return false, 0, nil
	}

	latestVersion := extractVersion(*latestRelease.TagName, resource.Resource.Source.TagFilter, versionRegex)
	if latestVersion == "" {
		u.logger.Infof("skip: %s, github release version unable to be parsed with expression '%s'", packageName, versionRegex)
		return false, 0, nil
	}

	currentVersion := ""
	if stateExists {
		state := &State{}
		err = u.readYAML(statePath, state)
		if err != nil {
			return false, 0, err
		}
		currentVersion = strings.TrimSpace(state.LatestVersion)
		if currentVersion == latestVersion {
			u.logger.Infof("skip: %s, on latest version", packageName)
			return false, 0, nil
		}
	}

	versionsToGenerate := []string{}
	if stateExists {
		allVersions, err := u.getAllVersions(owner, repository, resource.Resource.Source.TagFilter, versionRegex)
		if err != nil {
			return false, 0, err
		}
		if len(allVersions) == 0 {
			u.logger.Infof("skip: %s, no github releases found or unable to parse versions", packageName)
			return false, 0, nil
		}

		tempVersions := []string{}
		foundCurrent := false
		for _, version := range allVersions {
			if version == currentVersion {
				foundCurrent = true
				break
			}
			tempVersions = append(tempVersions, version)
		}

		if !foundCurrent {
			u.logger.Infof("warning: %s, current version '%s' not found in releases, generating all versions", packageName, currentVersion)
		}

		for i := len(tempVersions) - 1; i >= 0; i-- {
			versionsToGenerate = append(versionsToGenerate, tempVersions[i])
		}
	} else {
		versionsToGenerate = []string{latestVersion}
		if err := u.saveState(statePath, latestVersion); err != nil {
			return false, 0, err
		}
	}

	generated := 0
	for _, version := range versionsToGenerate {
		if strings.TrimSpace(version) == "" {
			continue
		}

		versionDir := u.path.Join(packageDir, version)
		err := u.fs.MkdirAll(versionDir, 0775)
		if err != nil {
			return false, generated, err
		}

		err = u.saveState(statePath, version)
		if err != nil {
			return false, generated, err
		}

		rendered, err := u.renderTemplate(packageDir)
		if err != nil {
			return false, generated, err
		}

		manifestPath := u.path.Join(versionDir, "package.yml")
		err = u.fs.WriteFile(manifestPath, rendered, 0644)
		if err != nil {
			return false, generated, err
		}
		generated++
	}

	err = u.saveState(statePath, latestVersion)
	if err != nil {
		return false, generated, err
	}

	return generated > 0, generated, nil
}

func (u *updatePackages) getAllVersions(owner, repository, tagFilter, versionRegex string) ([]string, error) {
	versions := []string{}
	page := 1
	for {
		releases, hasNext, err := u.github.ListReleases(owner, repository, page, u.perPage)
		if err != nil {
			return nil, err
		}
		if len(releases) == 0 {
			break
		}

		for _, release := range releases {
			if release == nil || release.TagName == nil {
				continue
			}
			version := extractVersion(*release.TagName, tagFilter, versionRegex)
			if strings.TrimSpace(version) == "" {
				continue
			}
			versions = append(versions, version)
		}

		if !hasNext {
			break
		}
		page++
	}
	return versions, nil
}

func (u *updatePackages) renderTemplate(packageDir string) ([]byte, error) {
	templatePath := u.path.Join(packageDir, "template.yml")
	templateContent, err := u.fs.ReadFile(templatePath)
	if err != nil {
		return nil, err
	}

	data := map[string]any{}
	for _, fileName := range []string{"platforms.yml", "resource.yml", "state.yml"} {
		settingsPath := u.path.Join(packageDir, fileName)
		ok, err := u.fs.Exists(settingsPath)
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}

		content, err := u.fs.ReadFile(settingsPath)
		if err != nil {
			return nil, err
		}

		parsed := map[string]any{}
		err = yaml.Unmarshal(content, &parsed)
		if err != nil {
			return nil, err
		}
		mergeMaps(data, parsed)
	}

	tmpl, err := template.New("template").Parse(string(templateContent))
	if err != nil {
		return nil, err
	}

	writer := &bytes.Buffer{}
	err = tmpl.Execute(writer, data)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func (u *updatePackages) saveState(statePath, version string) error {
	state := &State{LatestVersion: version}
	content, err := yaml.Marshal(state)
	if err != nil {
		return err
	}
	return u.fs.WriteFile(statePath, content, 0644)
}

func (u *updatePackages) readYAML(filePath string, out any) error {
	content, err := u.fs.ReadFile(filePath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, out)
}

func extractVersion(tag, tagFilter, versionRegex string) string {
	if strings.TrimSpace(tagFilter) != "" {
		tagFilterRegex, err := regexp.Compile(tagFilter)
		if err != nil {
			return ""
		}
		if !tagFilterRegex.MatchString(tag) {
			return ""
		}
	}

	regex, err := regexp.Compile(versionRegex)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(regex.FindString(tag))
}

func mergeMaps(target map[string]any, source map[string]any) {
	for key, sourceValue := range source {
		targetValue, ok := target[key]
		if !ok {
			target[key] = sourceValue
			continue
		}

		targetMap, targetOK := targetValue.(map[string]any)
		sourceMap, sourceOK := sourceValue.(map[string]any)
		if targetOK && sourceOK {
			mergeMaps(targetMap, sourceMap)
			continue
		}

		target[key] = sourceValue
	}
}

type defaultGitHubReleaseClient struct {
	client *github.Client
}

func (d *defaultGitHubReleaseClient) GetLatestRelease(owner, repository string) (*github.RepositoryRelease, error) {
	release, _, err := d.client.Repositories.GetLatestRelease(context.Background(), owner, repository)
	if err != nil {
		return nil, err
	}
	return release, nil
}

func (d *defaultGitHubReleaseClient) ListReleases(owner, repository string, page, perPage int) ([]*github.RepositoryRelease, bool, error) {
	releases, response, err := d.client.Repositories.ListReleases(context.Background(), owner, repository, &github.ListOptions{
		Page:    page,
		PerPage: perPage,
	})
	if err != nil {
		return nil, false, err
	}
	nextPage := response != nil && response.NextPage > 0
	return releases, nextPage, nil
}
