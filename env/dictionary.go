package env

import (
	"fmt"
	"os"

	"github.com/patrickhuber/wrangle/collections"
)

type dictionary struct {
}

// NewDictionary creates a dictionary wrapper around environment variables
func NewDictionary() collections.Dictionary {
	return &dictionary{}
}

func (d *dictionary) Get(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("unable to lookup environment variable '%s'", key)
	}
	return value, nil
}

func (d *dictionary) Set(key, value string) error {
	return os.Setenv(key, value)
}

func (d *dictionary) Lookup(key string) (string, bool) {
	return os.LookupEnv(key)
}

func (d *dictionary) Unset(key string) error {
	return os.Unsetenv(key)
}

func (d *dictionary) Keys() []string {
	return os.Environ()
}
