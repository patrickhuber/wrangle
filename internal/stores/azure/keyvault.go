package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/patrickhuber/wrangle/internal/stores"
)

type Store struct {
	uri string
}

type KeyVaultOptions struct {
	ClientID     string
	ClientSecret string
}

// NewKeyVault returns a new Store implemented by key vault
func NewKeyVault(uri string, options *KeyVaultOptions) *Store {
	return &Store{
		uri: uri,
	}
}

func (s Store) Get(key stores.Key) (any, error) {
	client, err := s.client()
	if err != nil {
		return nil, err
	}

	version := ""
	if !key.Data.Version.Latest {
		version = key.Data.Version.Value
	}

	resp, err := client.GetSecret(context.Background(), key.Data.Name, version, nil)
	if err != nil {
		return "", err
	}

	return *resp.Value, nil
}

func (s Store) Lookup(key stores.Key) (any, bool, error) {
	client, err := s.client()
	if err != nil {
		return nil, false, err
	}

	version := ""
	if !key.Data.Version.Latest {
		version = key.Data.Version.Value
	}

	resp := client.NewListSecretVersionsPager(key.Data.Name, nil)
	for resp.More() {
		page, err := resp.NextPage(context.Background())
		if err != nil {
			return nil, false, err
		}
		for _, secret := range page.SecretListResult.Value {
			if version != "" && version != secret.ID.Version() {
				continue
			}

			resp, err := client.GetSecret(
				context.Background(),
				key.Data.Name,
				version, nil)
			if err != nil {
				return nil, false, err
			}
			return *resp.Value, true, nil
		}
	}

	// unable to find secret
	return nil, false, nil
}

func (s Store) client() (*azsecrets.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return azsecrets.NewClient(s.uri, cred, nil)
}
