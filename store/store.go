package store

type ConfigStore interface {
	GetName() string
	Put(key string, value string) (string, error)
	GetByName(name string) (StoreData, error)
	GetByID(id string) (StoreData, error)
	Delete(key string) (int, error)
}
