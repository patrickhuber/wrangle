package keyring

import (
	"encoding/json"
	"errors"

	"github.com/99designs/keyring"
	"github.com/patrickhuber/wrangle/internal/stores"
)

func NewVault(service string) stores.Store {
	return &Vault{
		service: service,
	}
}

type Vault struct {
	service string
}

// Set sets the key to the vault in the key vault
func (v *Vault) Set(key stores.Key, value any) error {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: v.service,
	})
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return ring.Set(keyring.Item{
		Key:  key.Data.Name,
		Data: bytes,
	})
}

// Get returns the vault of the given key from the key vault
func (v *Vault) Get(k stores.Key) (any, bool, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: v.service,
	})
	if err != nil {
		return nil, false, err
	}

	item, err := ring.Get(k.Data.Name)
	if err != nil {
		if errors.Is(err, keyring.ErrKeyNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var value any
	err = json.Unmarshal(item.Data, &value)
	if err != nil {
		return nil, false, err
	}

	return value, true, nil
}

// List implements stores.Store.
func (v *Vault) List() ([]stores.Key, error) {

	ring, err := keyring.Open(keyring.Config{
		ServiceName: v.service,
	})
	if err != nil {
		return nil, err
	}

	keys, err := ring.Keys()
	if err != nil {
		return nil, err
	}

	var result []stores.Key
	for _, key := range keys {
		result = append(result, stores.Key{
			Data: stores.Data{
				Name: key,
			},
		})
	}
	return result, nil
}
