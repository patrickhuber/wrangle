package store

// Store represents a data store for config
type Store interface {
	GetName() string
	GetType() string
	Put(key string, value string) (string, error)
	GetByName(name string) (Data, error)
	Delete(key string) (int, error)
}
