package memory

import (
	"fmt"

	"github.com/patrickhuber/wrangle/store"
)

// MemoryStore - Struct that represents a memory store
type memoryStore struct {
	name string
	data map[string]store.Data
}

// NewMemoryStore - Creates a new memory store with the given name
func NewMemoryStore(name string) store.Store {
	data := map[string]store.Data{}
	return &memoryStore{
		name: name,
		data: data,
	}
}

// Name - Gets the name for the memory store
func (s *memoryStore) Name() string {
	return s.name
}

// Type - Gets the type for the store. Always "memory"
func (s *memoryStore) Type() string {
	return "memory"
}

// Set - Puts the config value under the value in the memory store
func (s *memoryStore) Set(key string, value string) (string, error) {
	data := store.NewData(
		key,
		value,
	)
	s.data[key] = data
	return value, nil
}

// Get - Gets the config value by name
func (s *memoryStore) Get(key string) (store.Data, error) {
	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("unable to locate key %s", key)
	}
	return value, nil
}

// Delete - Deletes the value from the config store
func (s *memoryStore) Delete(key string) error {
	data, err := s.Get(key)
	if err != nil {
		return err
	}
	delete(s.data, data.Name())
	return nil
}
