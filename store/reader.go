package store

// Reader defines a read only interface for interacting with stores
type Reader interface {
	Get(key string) (Item, error)
	List(path string) ([]Item, error)
}
