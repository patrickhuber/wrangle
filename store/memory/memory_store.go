package memory

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/store"
)

// MemoryStore - Struct that represents a memory store
type memoryStore struct {
	name string
	data map[string]store.Item
}

// NewMemoryStore - Creates a new memory store with the given name
func NewMemoryStore(name string) store.Store {
	data := map[string]store.Item{}
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
func (s *memoryStore) Set(item store.Item) error {
	key := s.normalizeKey(item.Name())
	s.data[key] = item
	return nil
}

// Get - Gets the config value by name
func (s *memoryStore) Get(key string) (store.Item, error) {
	key = s.normalizeKey(key)
	value, ok := s.data[key]
	if !ok {
		return nil, fmt.Errorf("unable to locate key %s", key)
	}
	return value, nil
}

func (s *memoryStore) List(path string) ([]store.Item, error) {
	path = s.normalizeKey(path)
	pathSplit := strings.Split(path, "/")
	items := []store.Item{}
	for k, v := range s.data {
		if len(pathSplit) == 2 && pathSplit[1] == "" {
			items = append(items, v)
			continue
		}

		keySplit := strings.Split(k, "/")
		if len(pathSplit) > len(keySplit) {
			continue
		}
		isMatch := true
		for i := 1; i < len(pathSplit); i++ {
			if pathSplit[i] != keySplit[i] {
				isMatch = false
				break
			}
		}
		if isMatch {
			items = append(items, v)
		}
	}
	return items, nil
}

// Delete - Deletes the value from the config store
func (s *memoryStore) Delete(key string) error {
	key = s.normalizeKey(key)
	_, err := s.Get(key)
	if err != nil {
		return err
	}
	delete(s.data, key)
	return nil
}

func (s *memoryStore) normalizeKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		return "/" + key
	}
	return key
}
