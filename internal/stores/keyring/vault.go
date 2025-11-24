package keyring

import (
	"encoding/json"
	"errors"

	"github.com/99designs/keyring"
	"github.com/patrickhuber/wrangle/internal/stores"
)

func NewVault(config keyring.Config) stores.Store {
	return &Vault{
		config: config,
	}
}

type Vault struct {
	config keyring.Config
}

func (v *Vault) open() (keyring.Keyring, error) {
	return keyring.Open(v.config)
}

// Set sets the key to the vault in the key vault
func (v *Vault) Set(key stores.Key, value any) error {

	ring, err := v.open()
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

	ring, err := v.open()
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

	ring, err := v.open()
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
				Version: stores.Version{
					Latest: true,
				},
			},
		})
	}
	return result, nil
}
