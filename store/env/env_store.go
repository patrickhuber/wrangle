package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/patrickhuber/wrangle/store"
)

type envStore struct {
	name   string
	lookup map[string]string
}

func NewEnvStore(name string, lookup map[string]string) store.Store {
	return &envStore{
		name:   name,
		lookup: lookup}
}

func (s *envStore) Name() string {
	return s.name
}

func (s *envStore) Type() string {
	return "env"
}

func (s *envStore) GetByName(name string) (store.Data, error) {
	// cleanup the key just in case there is a forward slash
	name = s.cleanKey(name)

	// lookup the variable
	environmentVariableName, ok := s.lookup[name]
	if !ok {
		return nil, fmt.Errorf("name '%s' not found in definition store '%s'", name, s.Name())
	}

	// look up the environment variable
	data := os.Getenv(environmentVariableName)

	return store.NewData(name, name, data), nil
}

func (s *envStore) Delete(name string) (int, error) {
	return 1, fmt.Errorf("Delete method not implemented")
}

func (s *envStore) Put(key string, value string) (string, error) {
	return "", fmt.Errorf("Put method not implemented")
}

func (s *envStore) cleanKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		return key
	}
	return key[1:]
}
