package local

import (
	"github.com/99designs/keyring"
	"github.com/patrickhuber/wrangle/store"
)

type localStore struct {
}

func (s *localStore) Get(key string) (store.Item, error) {

	ring, err := loadKeyRing()
	if err != nil {
		return nil, err
	}

	item, err := ring.Get(key)
	if err != nil {
		return nil, err
	}

	return store.NewData(key, item.Data), nil
}

func (s *localStore) Set(key string, value string) (string, error) {
	ring, err := loadKeyRing()
	if err != nil {
		return "", err
	}

	item := keyring.Item{
		Key:  key,
		Data: []byte(value),
	}

	err = ring.Set(item)
	if err != nil {
		return "", err
	}

	return value, nil
}

func (s *localStore) Delete(key string) error {
	ring, err := loadKeyRing()
	if err != nil {
		return err
	}
	return ring.Remove(key)
}

func loadKeyRing() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: "wrangle",
	})
}
