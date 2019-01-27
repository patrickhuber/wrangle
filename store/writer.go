package store

// Writer defines a write only store
type Writer interface {
	Set(item Item) error
	Delete(key string) error
}
