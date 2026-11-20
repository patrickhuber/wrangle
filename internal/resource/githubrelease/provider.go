package githubrelease

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/patrickhuber/go-log"
	"github.com/patrickhuber/wrangle/internal/resource"
)

type ClientFactory func(token string) GitHub

type provider struct {
	clientFactory ClientFactory
	logger        log.Logger
}

func NewProvider(logger log.Logger) resource.Provider {
	return NewProviderWithClientFactory(NewGitHubClient, logger)
}

func NewProviderWithClientFactory(factory ClientFactory, logger log.Logger) resource.Provider {
	return &provider{clientFactory: factory, logger: logger}
}

func (*provider) Type() string {
	return "github-release"
}

func (p *provider) Create(config resource.Config) (resource.Resource, error) {
	source, err := sourceFromConfig(config.Source)
	if err != nil {
		return nil, err
	}
	client := p.clientFactory(source.AccessToken)
	return NewResource(config.Name, source, NewService(client, p.logger)), nil
}

func sourceFromConfig(values map[string]string) (Source, error) {
	allowedKeys := map[string]bool{
		"owner":         true,
		"repository":    true,
		"version_regex": true,
		"version-regex": true,
		"tag_filter":    true,
		"access_token":  true,
	}
	for _, key := range sortedSourceKeys(values) {
		if !allowedKeys[key] {
			return Source{}, fmt.Errorf("unknown github-release source key %q", key)
		}
	}

	owner := strings.TrimSpace(values["owner"])
	if owner == "" {
		return Source{}, fmt.Errorf("github-release source requires %q", "owner")
	}
	repository := strings.TrimSpace(values["repository"])
	if repository == "" {
		return Source{}, fmt.Errorf("github-release source requires %q", "repository")
	}

	versionRegex := values["version_regex"]
	if versionRegex == "" {
		versionRegex = values["version-regex"]
	}
	tagFilter := values["tag_filter"]
	for _, entry := range []struct {
		key     string
		pattern string
	}{
		{key: "version_regex", pattern: versionRegex},
		{key: "tag_filter", pattern: tagFilter},
	} {
		if entry.pattern == "" {
			continue
		}
		if _, err := regexp.Compile(entry.pattern); err != nil {
			return Source{}, fmt.Errorf("invalid github-release source %q: %w", entry.key, err)
		}
	}

	return Source{
		Owner:        owner,
		Repository:   repository,
		VersionRegex: versionRegex,
		TagFilter:    tagFilter,
		AccessToken:  values["access_token"],
	}, nil
}

func sortedSourceKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
