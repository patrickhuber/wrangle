package memory

import (
	"fmt"

	"github.com/patrickhuber/wrangle/store"

	uuid "github.com/google/uuid"
)

// MemoryStore - Struct that represents a memory store
type memoryStore struct {
	name    string
	data    map[string]store.Data
	keyToID map[string]string
}

// NewMemoryStore - Creates a new memory store with the given name
func NewMemoryStore(name string) store.Store {
	data := map[string]store.Data{}
	keyToID := map[string]string{}
	return &memoryStore{
		name:    name,
		data:    data,
		keyToID: keyToID,
	}
}

// Name - Gets the name for the memory store
func (store *memoryStore) Name() string {
	return store.name
}

// Type - Gets the type for the store. Always "memory"
func (store *memoryStore) Type() string {
	return "memory"
}

// Put - Puts the config value under the value in the memory store
func (memoryStore *memoryStore) Put(key string, value string) (string, error) {
	data := store.NewData(
		uuid.New().String(),
		key,
		value,
	)
	memoryStore.data[data.ID()] = data
	memoryStore.keyToID[key] = data.ID()
	return data.ID(), nil
}

// GetByName - Gets the config value by name
func (memoryStore *memoryStore) GetByName(key string) (store.Data, error) {
	id, ok := memoryStore.keyToID[key]
	if !ok {
		return nil, fmt.Errorf("Unable to find key %s", key)
	}
	return memoryStore.GetByID(id)
}

// GetByID - Gets the value by ID
func (memoryStore *memoryStore) GetByID(id string) (store.Data, error) {
	value, ok := memoryStore.data[id]
	if ok != true {
		return nil, fmt.Errorf("Unable to find id %s", id)
	}
	return value, nil
}

// Delete - Deletes the value from the config store
func (memoryStore *memoryStore) Delete(key string) (int, error) {
	data, err := memoryStore.GetByName(key)
	if err != nil {
		return 0, err
	}
	delete(memoryStore.keyToID, key)
	delete(memoryStore.data, data.ID())
	return 1, nil
}
