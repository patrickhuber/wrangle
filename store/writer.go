package store

// Writer defines a write only store
type Writer interface {
	Set(key string, value string) (string, error)
	Delete(key string) (int, error)
}
