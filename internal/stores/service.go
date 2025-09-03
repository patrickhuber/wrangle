package stores

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/config"
)

type Service interface {
	Get(name string) (Store, error)
	List() ([]Store, error)
}

type service struct {
	configuration config.Configuration
	registry      Registry
	cache         map[string]Store
}

func NewService(configuration config.Configuration, registry Registry) Service {
	return &service{
		configuration: configuration,
		registry:      registry,
		cache:         map[string]Store{},
	}
}

func (s *service) Get(name string) (Store, error) {
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

func (s *service) List() ([]Store, error) {

	err := s.load()
	if err != nil {
		return nil, err
	}

	var results []Store
	for _, v := range s.cache {
		results = append(results, v)
	}
	return results, nil
}

func (s *service) load() error {
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
