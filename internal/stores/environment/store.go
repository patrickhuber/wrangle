package environment

import (
	"fmt"
	"path"

	"github.com/patrickhuber/go-cross/env"
	"github.com/patrickhuber/wrangle/internal/stores"
)

type Store struct {
	environment env.Environment
	include     []string
	includeSet  bool
	exclude     []string
}

func (s *Store) Get(key stores.Key) (any, bool, error) {
	name, err := variableName(key)
	if err != nil {
		return nil, false, err
	}
	included, err := s.includes(name)
	if err != nil {
		return nil, false, err
	}
	if !included {
		return nil, false, nil
	}
	value, found := s.environment.Lookup(name)
	if !found {
		return nil, false, nil
	}
	return value, true, nil
}

func (s *Store) Set(_ stores.Key, _ any) error {
	return fmt.Errorf("environment store is read-only")
}

func (s *Store) List() ([]stores.Key, error) {
	var keys []stores.Key
	for name := range s.environment.Export() {
		included, err := s.includes(name)
		if err != nil {
			return nil, err
		}
		if !included {
			continue
		}
		keys = append(keys, stores.Key{Data: stores.Data{Name: name, Version: stores.Version{Latest: true}}})
	}
	return keys, nil
}

func (s *Store) includes(name string) (bool, error) {
	if s.includeSet {
		matched, err := matchesAny(name, s.include)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	matched, err := matchesAny(name, s.exclude)
	if err != nil {
		return false, err
	}
	return !matched, nil
}

func matchesAny(name string, patterns []string) (bool, error) {
	for _, pattern := range patterns {
		matched, err := path.Match(pattern, name)
		if err != nil {
			return false, err
		}
		if matched {
			return true, nil
		}
	}
	return false, nil
}

func variableName(key stores.Key) (string, error) {
	if len(key.Path.Segments) != 0 || !key.Data.Version.Latest {
		return "", fmt.Errorf("environment variable key cannot include a path or version")
	}
	if key.Data.Name == "" {
		return "", fmt.Errorf("environment variable name is required")
	}
	return key.Data.Name, nil
}
