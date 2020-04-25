package feed

import (
	"strings"
)

// PackageVersionManifest defines a pacakge version manifest item and its content
type PackageVersionManifest struct {
	Content string
	Name    string
}

// PackageVersion defines a package version response object
type PackageVersion struct {
	Version  string
	Manifest *PackageVersionManifest
	Feeds    []string
}

// Package defines a package response object
type Package struct {
	Name     string
	Versions []*PackageVersion
	Latest   string
}

// ListRequest contains the request for listing packages
type ListRequest struct {
}

// ListResponse contains the response from listing packages
type ListResponse struct {
	Packages []*Package
}

// GetRequest defines criteria for fetching a specific package
type GetRequest struct {
	Name           string
	Version        string
	IncludeContent bool
}

// GetResponse Contains the list of matching packages
type GetResponse struct {
	Package *Package
}

// CreateRequest contains the list of packages to create
type CreateRequest struct {
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

// CreateResponse creates a package and returns the created ID
type CreateResponse struct {
	Packages []*Package
}

// LatestRequest defines a query for the lastest version of a package
type LatestRequest struct {
	Name string
}

// LatestResponse defines a response for a FeedLastestRequest query
type LatestResponse struct {
	Package *Package
}

// InfoRequest requests feed information
type InfoRequest struct {
}

// InfoResponse returns information for the feed
type InfoResponse struct {
	URI string
}

// Service defines a package feed service
type Service interface {
	List(request *ListRequest) (*ListResponse, error)
	Get(request *GetRequest) (*GetResponse, error)
	Create(request *CreateRequest) (*CreateResponse, error)
	Latest(request *LatestRequest) (*LatestResponse, error)
	Info(request *InfoRequest) (*InfoResponse, error)
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

type packageApply struct {
	Latest bool
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
