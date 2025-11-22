package vault

import (
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/patrickhuber/wrangle/internal/stores"
)

type Store struct {
	client *api.Client
	path   string
}

type VaultOptions struct {
	Address string
	Token   string
	Path    string
}

// NewVault returns a new Store implemented by HashiCorp Vault
func NewVault(options *VaultOptions) (*Store, error) {
	config := api.DefaultConfig()
	if options.Address != "" {
		config.Address = options.Address
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}

	if options.Token != "" {
		client.SetToken(options.Token)
	}

	path := "secret"
	if options.Path != "" {
		path = options.Path
	}

	return &Store{
		client: client,
		path:   path,
	}, nil
}

func (s *Store) Get(key stores.Key) (any, bool, error) {
	secretPath := fmt.Sprintf("%s/data/%s", s.path, key.Data.Name)

	secret, err := s.client.Logical().Read(secretPath)
	if err != nil {
		return nil, false, err
	}

	if secret == nil {
		return nil, false, nil
	}

	// KV v2 stores data under "data" key
	data, ok := secret.Data["data"]
	if !ok {
		return nil, false, nil
	}

	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return nil, false, fmt.Errorf("unexpected data format in vault secret")
	}

	// Get the value from the data map
	value, ok := dataMap["value"]
	if !ok {
		// If no "value" key, return the entire data map
		return dataMap, true, nil
	}

	return value, true, nil
}

func (s *Store) Set(key stores.Key, value any) error {
	secretPath := fmt.Sprintf("%s/data/%s", s.path, key.Data.Name)

	// Wrap value in data map for KV v2
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"value": value,
		},
	}

	_, err := s.client.Logical().Write(secretPath, data)
	return err
}

func (s *Store) List() ([]stores.Key, error) {
	listPath := fmt.Sprintf("%s/metadata", s.path)

	secret, err := s.client.Logical().List(listPath)
	if err != nil {
		return nil, err
	}

	if secret == nil {
		return []stores.Key{}, nil
	}

	keysData, ok := secret.Data["keys"]
	if !ok {
		return []stores.Key{}, nil
	}

	keysList, ok := keysData.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected keys format in vault list response")
	}

	var keys []stores.Key
	for _, k := range keysList {
		keyName, ok := k.(string)
		if !ok {
			continue
		}

		keys = append(keys, stores.Key{
			Data: stores.Data{
				Name: keyName,
			},
		})
	}

	return keys, nil
}
