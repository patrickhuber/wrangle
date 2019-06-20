package store

// Reader defines a read only interface for interacting with stores
type Reader interface {
	Get(key string) (Item, error)
	Lookup(key string)(Item, bool, error)
	List(path string) ([]Item, error)
}
