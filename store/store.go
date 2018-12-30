package store

// Store represents a data store for config
type Store interface {
	Name() string
	Type() string
	Reader
	Writer
}
