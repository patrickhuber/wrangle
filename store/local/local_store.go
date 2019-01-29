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

	return store.NewItem(key, store.Value, item.Data), nil
}

func (s *localStore) Set(item store.Item) error {
	ring, err := loadKeyRing()
	if err != nil {
		return err
	}

	key := item.Name()
	// need to use structure -> byte array serialization here
	value := []byte{}
	keyringItem := keyring.Item{
		Key:  key,
		Data: value,
	}

	err = ring.Set(keyringItem)
	if err != nil {
		return err
	}

	return nil
}

func (s *localStore) Delete(key string) error {
	ring, err := loadKeyRing()
	if err != nil {
		return err
	}
	return ring.Remove(key)
}

func (s *localStore) Copy(item store.Item, destination string) error {
	return nil
}

func loadKeyRing() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: "wrangle",
	})
}
