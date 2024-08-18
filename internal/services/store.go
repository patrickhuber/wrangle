package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

type Store interface {
	Get(name string) (stores.Store, error)
	List() ([]stores.Store, error)
}

type store struct {
	configuration Configuration
	registry      stores.Registry
	cache         map[string]stores.Store
}

func NewStore(configuration Configuration, registry stores.Registry) Store {
	return &store{
		configuration: configuration,
		registry:      registry,
		cache:         map[string]stores.Store{},
	}
}

func (s *store) Get(name string) (stores.Store, error) {
	err := s.load()
	if err != nil {
		return nil, err
	}

	st, ok := s.cache[name]
	if !ok {
		return nil, fmt.Errorf("unable to locate store '%s' in configuration", name)
	}
	return st, nil
}

func (s *store) List() ([]stores.Store, error) {

	err := s.load()
	if err != nil {
		return nil, err
	}

	var results []stores.Store
	for _, v := range s.cache {
		results = append(results, v)
	}
	return results, nil
}

func (s *store) load() error {
	if len(s.cache) > 0 {
		return nil
	}
	cfg, err := s.configuration.Get()
	if err != nil {
		return err
	}
	for _, c := range cfg.Spec.Stores {
		factory, err := s.registry.Get(c.Type)
		if err != nil {
			return err
		}
		s.cache[c.Name], err = factory.Create(c.Properties)
		if err != nil {
			return err
		}
	}
	return nil
}
