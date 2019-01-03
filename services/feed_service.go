package services

// FeedListRequest contains the request for listing packages
type FeedListRequest struct {
}

// FeedListResponsePackage defines a package
type FeedListResponsePackage struct {
	Name     string
	Versions []string
}

// FeedListResponse contains the response from listing packages
type FeedListResponse struct {
	Packages []FeedListResponsePackage
}

// FeedService defines a package feed service
type FeedService interface {
	List(request *FeedListRequest) (*FeedListResponse, error)
}
