package vault

import (
	"context"
	"fmt"

	"github.com/hashicorp/vault/api"
	"github.com/patrickhuber/wrangle/internal/stores"
)

type Store struct {
	client *api.Client
	kv     *api.KVv2
	path   string
}

type VaultOptions struct {
	Address  string
	Token    string
	Path     string
	RoleID   string
	SecretID string
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

	// AppRole authentication takes precedence over token
	if options.RoleID != "" && options.SecretID != "" {
		err = authenticateWithAppRole(client, options.RoleID, options.SecretID)
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate with AppRole: %w", err)
		}
	} else if options.Token != "" {
		client.SetToken(options.Token)
	}

	path := "secret"
	if options.Path != "" {
		path = options.Path
	}

	// Use the KVv2 API for KV version 2 secrets engine
	kv := client.KVv2(path)

	return &Store{
		client: client,
		kv:     kv,
		path:   path,
	}, nil
}

// authenticateWithAppRole authenticates to Vault using AppRole auth method
func authenticateWithAppRole(client *api.Client, roleID, secretID string) error {
	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	resp, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		return err
	}

	if resp == nil || resp.Auth == nil {
		return fmt.Errorf("authentication response is nil")
	}

	client.SetToken(resp.Auth.ClientToken)
	return nil
}

// Client returns the underlying Vault API client for advanced operations.
// This method is primarily intended for testing purposes to configure Vault
// (e.g., setting up AppRole authentication). Production code should use the
// Store interface methods (Get, Set, List) instead.
func (s *Store) Client() *api.Client {
	return s.client
}

func (s *Store) Get(key stores.Key) (any, bool, error) {
	ctx := context.Background()
	
	secret, err := s.kv.Get(ctx, key.Data.Name)
	if err != nil {
		// Check if it's a 404 error (secret not found)
		if respErr, ok := err.(*api.ResponseError); ok && respErr.StatusCode == 404 {
			return nil, false, nil
		}
		return nil, false, err
	}

	if secret == nil || secret.Data == nil {
		return nil, false, nil
	}

	// Get the value from the data map
	value, ok := secret.Data["value"]
	if !ok {
		// If no "value" key, return the entire data map
		return secret.Data, true, nil
	}

	return value, true, nil
}

func (s *Store) Set(key stores.Key, value any) error {
	ctx := context.Background()
	
	// Create data map with the value
	data := map[string]interface{}{
		"value": value,
	}

	_, err := s.kv.Put(ctx, key.Data.Name, data)
	return err
}

func (s *Store) List() ([]stores.Key, error) {
	// For listing, we still need to use the Logical interface
	// The KVv2 API doesn't expose a direct List method
	listPath := fmt.Sprintf("%s/metadata", s.path)

	secret, err := s.client.Logical().List(listPath)
	if err != nil {
		return nil, err
	}

	if secret == nil || secret.Data == nil {
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
