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

func (s *envStore) Get(key string) (store.Data, error) {
	// cleanup the key just in case there is a forward slash
	key = s.cleanKey(key)

	// lookup the variable
	environmentVariableName, ok := s.lookup[key]
	if !ok {
		return nil, fmt.Errorf("key '%s' not found in definition store '%s'", key, s.Name())
	}

	// look up the environment variable
	data, ok := os.LookupEnv(environmentVariableName)
	if !ok {
		return nil, fmt.Errorf("variable '%s' is not set in the environment variables", environmentVariableName)
	}

	return store.NewData(key, data), nil
}

func (s *envStore) Delete(key string) error {
	return fmt.Errorf("Delete method not implemented")
}

func (s *envStore) Set(key string, value string) (string, error) {
	return "", fmt.Errorf("Set method not implemented")
}

func (s *envStore) cleanKey(key string) string {
	if !strings.HasPrefix(key, "/") {
		return key
	}
	return key[1:]
}
