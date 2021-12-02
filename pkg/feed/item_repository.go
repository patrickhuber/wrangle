package feed

// ItemRepository provides a repository for items
type ItemRepository interface {
	Get(name string) (*Item, error)
	List(where []*ItemReadAnyOf) ([]*Item, error)
}
