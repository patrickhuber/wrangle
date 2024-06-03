package services

import (
	"fmt"

	"github.com/patrickhuber/wrangle/internal/stores"
)

type Secret struct {
	store Store
}

func NewSecret(st Store) Secret {
	return Secret{
		store: st,
	}
}

func (s *Secret) Get(store string, key string) (any, error) {

	st, err := s.store.Get(store)
	if err != nil {
		return nil, err
	}

	v, ok, err := st.Get(stores.Key{Data: stores.Data{Name: key}})
	if err != nil {
		return "", err
	}
	if !ok {
		return "", fmt.Errorf("unable to locate key %s in store %s", key, store)
	}

	return v, nil
}

func (s *Secret) Set(store string, key string, value string) error {
	st, err := s.store.Get(store)
	if err != nil {
		return err
	}
	return st.Set(stores.Key{Data: stores.Data{Name: key}}, value)
}
