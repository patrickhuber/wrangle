package env

import (
	"fmt"
	"strings"

	"github.com/patrickhuber/wrangle/collections"

	"github.com/patrickhuber/wrangle/store"
)

type envStore struct {
	name      string
	lookup    map[string]string
	variables collections.Dictionary
}

func NewEnvStore(
	name string,
	lookup map[string]string,
	variables collections.Dictionary) store.Store {
	return &envStore{
		name:      name,
		lookup:    lookup,
		variables: variables}
}

func (s *envStore) Name() string {
	return s.name
}

func (s *envStore) Type() string {
	return "env"
}

func (s *envStore) Get(key string) (store.Item, error) {
	// cleanup the key just in case there is a forward slash
	key = s.cleanKey(key)

	// lookup the variable
	environmentVariableName, ok := s.lookup[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in definition store '%s'", key, s.Name())
	}

	// look up the environment variable
	data, ok := s.variables.Lookup(environmentVariableName)
	if !ok {
		return nil, fmt.Errorf("variable '%s' is not set in the environment variables", environmentVariableName)
	}

	return store.NewValueItem(key, data), nil
}

func (s *envStore) List(path string) ([]store.Item, error) {
	return nil, nil
}

func (s *envStore) Lookup(path string) (store.Item, bool, error){
	// cleanup the key just in case there is a forward slash
	key := s.cleanKey(path)

	// lookup the variable
	environmentVariableName, ok := s.lookup[key]
	if !ok {
		return nil, false, nil
	}

	// look up the environment variable
	data, ok := s.variables.Lookup(environmentVariableName)
	if !ok {
		return nil, false, nil
	}
	
	return store.NewValueItem(key, data), true, nil
}

func (s *envStore) Delete(key string) error {
	return fmt.Errorf("Delete method not implemented")
}

func (s *envStore) Set(item store.Item) error {
	return fmt.Errorf("Set method not implemented")
}

func (s *envStore) Copy(item store.Item, destination string) error {
	return nil
}

func (s *envStore) cleanKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		return key
	}
	return key[1:]
}
