package memory

import (
	"fmt"

	"github.com/patrickhuber/cli-mgr/store"

	uuid "github.com/google/uuid"
)

// MemoryStore - Struct that represents a memory store
type MemoryStore struct {
	Name    string
	Data    map[string]store.Data
	KeyToID map[string]string
}

// NewMemoryStore - Creates a new memory store with the given name
func NewMemoryStore(name string) *MemoryStore {
	data := map[string]store.Data{}
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

// GetType - Gets the type for the store. Always "memory"
func (store *MemoryStore) GetType() string {
	return "memory"
}

// Put - Puts the config value under the value in the memory store
func (memoryStore *MemoryStore) Put(key string, value string) (string, error) {
	data := store.NewData(
		uuid.New().String(),
		key,
		value,
	)
	memoryStore.Data[data.GetID()] = data
	memoryStore.KeyToID[key] = data.GetID()
	return data.GetID(), nil
}

// GetByName - Gets the config value by name
func (memoryStore *MemoryStore) GetByName(key string) (store.Data, error) {
	id, ok := memoryStore.KeyToID[key]
	if !ok {
		return nil, fmt.Errorf("Unable to find key %s", key)
	}
	return memoryStore.GetByID(id)
}

// GetByID - Gets the value by ID
func (memoryStore *MemoryStore) GetByID(id string) (store.Data, error) {
	value, ok := memoryStore.Data[id]
	if ok != true {
		return nil, fmt.Errorf("Unable to find id %s", id)
	}
	return value, nil
}

// Delete - Deletes the value from the config store
func (memoryStore *MemoryStore) Delete(key string) (int, error) {
	data, err := memoryStore.GetByName(key)
	if err != nil {
		return 0, err
	}
	delete(memoryStore.KeyToID, key)
	delete(memoryStore.Data, data.GetID())
	return 1, nil
}
