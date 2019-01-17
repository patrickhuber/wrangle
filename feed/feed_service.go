package feed

import "strings"

type PackageVersionManifest struct {
	Content *string
	Name    string
}

type PackageVersion struct {
	Version  string
	Manifest *PackageVersionManifest
	Feeds    []string
}

type Package struct {
	Name     string
	Versions []*PackageVersion
}

// FeedListRequest contains the request for listing packages
type FeedListRequest struct {
}

// FeedListResponse contains the response from listing packages
type FeedListResponse struct {
	Packages []*Package
}

// FeedGetRequest defines criteria for fetching a specific package
type FeedGetRequest struct {
	Name           string
	Version        string
	IncludeContent bool
}

// FeedGetResponse Contains the list of matching packages
type FeedGetResponse struct {
	Package *Package
}

// FeedCreateRequest contains the list of packages to create
type FeedCreateRequest struct {
	Packages []*PackageCreate
}

type PackageCreate struct {
	Name     string
	Versions []*PackageVersionCreate
}

type PackageVersionCreate struct {
	Version  string
	Contents string
}

// FeedCreateResponse creates a package and returns the created ID
type FeedCreateResponse struct {
	Packages []*Package
}

// FeedService defines a package feed service
type FeedService interface {
	List(request *FeedListRequest) (*FeedListResponse, error)
	Get(request *FeedGetRequest) (*FeedGetResponse, error)
	Create(request *FeedCreateRequest) (*FeedCreateResponse, error)
}

type packageCriteria struct {
	Name     string
	Versions []string
}

type packageCriteriaWhere struct {
	Or []*packageCriteriaAnd
}

type packageCriteriaAnd struct {
	And []*packageCriteria
}

type packageInclude struct {
	Content bool
}

func evaluate(filter *packageCriteriaWhere, packageName string, packageVersion string) bool {
	if filter == nil || filter.Or == nil || len(filter.Or) == 0 {
		return true
	}
	for _, or := range filter.Or {
		if evaluateOr(or, packageName, packageVersion) {
			return true
		}
	}
	return false
}

func evaluateOr(or *packageCriteriaAnd, packageName string, packageVersion string) bool {
	if or == nil || or.And == nil || len(or.And) == 0 {
		return false
	}
	for _, and := range or.And {
		if !matchName(and.Name, packageName) || !matchAnyVersion(and.Versions, packageVersion) {
			return false
		}
	}
	return true
}

func matchName(name, packageName string) bool {
	// no filter was specified, return true
	if strings.TrimSpace(name) == "" {
		return true
	}
	return name == packageName
}

func matchAnyVersion(versions []string, version string) bool {
	// no filter was specified, return true
	if versions == nil || len(versions) == 0 {
		return true
	}

	for _, v := range versions {
		if v == version {
			return true
		}
	}

	return false
}
