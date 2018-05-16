package store

import (
	"fmt"

	uuid "github.com/google/uuid"
)

// MemoryStore - Struct that represents a memory store
type MemoryStore struct {
	Name    string
	Data    map[string]StoreData
	KeyToID map[string]string
}

// NewMemoryStore - Creates a new memory store with the given name
func NewMemoryStore(name string) *MemoryStore {
	data := map[string]StoreData{}
	keyToID := map[string]string{}
	return &MemoryStore{
		Name:    name,
		Data:    data,
		KeyToID: keyToID,
	}
}

// GetName - Gets the name for the memory store
func (store *MemoryStore) GetName() string {
	return store.Name
}

// Put - Puts the config value under the value in the memory store
func (store *MemoryStore) Put(key string, value string) (string, error) {
	data := StoreData{
		ID:    uuid.New().String(),
		Name:  key,
		Value: value,
	}
	store.Data[data.ID] = data
	store.KeyToID[key] = data.ID
	return data.ID, nil
}

// GetByName - Gets the config value by name
func (store *MemoryStore) GetByName(key string) (StoreData, error) {
	id, ok := store.KeyToID[key]
	if !ok {
		return StoreData{}, fmt.Errorf("Unable to find key %s", key)
	}
	return store.GetByID(id)
}

// GetByID - Gets the value by ID
func (store *MemoryStore) GetByID(id string) (StoreData, error) {
	value, ok := store.Data[id]
	if ok != true {
		return StoreData{}, fmt.Errorf("Unable to find id %s", id)
	}
	return value, nil
}

// Delete - Deletes the value from the config store
func (store *MemoryStore) Delete(key string) (int, error) {
	data, err := store.GetByName(key)
	if err != nil {
		return 0, err
	}
	delete(store.KeyToID, key)
	delete(store.Data, data.ID)
	return 1, nil
}
