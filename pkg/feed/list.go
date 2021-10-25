package feed

type ListRequest struct{}
type ListResponse struct {
	Items []*Item
}
