package githubrelease

import (
	"regexp"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// Service defines a github release service
type Service interface {
	// Check Returns a list of all versions that are greater than or equal to the current release
	Check(request *CheckRequest) (*CheckResponse, error)
}

type service struct {
	client GitHub
}

// NewService creates a new githubrelease service
func NewService(client GitHub) Service {
	return &service{
		client: client,
	}
}

func (s *service) Check(request *CheckRequest) (*CheckResponse, error) {
	listRequest := &ListRequest{
		Owner:      request.Source.Owner,
		Repository: request.Source.Repository,
	}
	githubReleases, err := s.client.ListReleases(listRequest)
	if err != nil {
		return nil, err
	}

	versions := []*Version{}

	// check for empty results, return empty list of versions
	if len(githubReleases) == 0 {
		return &CheckResponse{
			Versions: versions,
		}, nil
	}

	requestVersion, err := s.getVersion(request.Version.ID)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(request.Source.TagFilter)
	if err != nil {
		return nil, err
	}

	// iterate over the releases, convert to versions
	// once converted, check if it is greater than the request version
	for _, release := range githubReleases {
		if release.TagName == nil {
			continue
		}

		releaseVersionString := *release.TagName
		if !re.MatchString(releaseVersionString) {
			continue
		}

		v, err := semver.NewVersion(releaseVersionString)
		if err != nil {
			continue
		}

		// current version >= new version
		if v.LessThan(requestVersion) {
			continue
		}

		version := &Version{ID: v.String()}
		versions = append(versions, version)
	}

	return &CheckResponse{
		Versions: versions,
	}, nil
}

func (s *service) getVersion(version string) (*semver.Version, error) {

	var semverVersion *semver.Version
	var err error

	// if current version is empty, initialize to zero
	if len(strings.TrimSpace(version)) == 0 {
		semverVersion, err = semver.NewVersion("0.0.0")
		return semverVersion, err
	}

	// otherwise parse the current version
	semverVersion, err = semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	return semverVersion, nil
}
