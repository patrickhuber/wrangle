package store

type ConfigStore interface {
	GetName() string
	GetType() string
	Put(key string, value string) (string, error)
	GetByName(name string) (StoreData, error)
	Delete(key string) (int, error)
}
