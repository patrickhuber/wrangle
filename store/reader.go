package store

// Reader defines a read only interface for interacting with stores
type Reader interface {
	Get(key string) (Data, error)
}
