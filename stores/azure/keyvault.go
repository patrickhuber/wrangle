package azure

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

type Store struct {
	uri string
}

type KeyVaultOptions struct {
	ClientID     string
	ClientSecret string
}

// NewKeyVault returns a new Store implemented by key vault
func NewKeyVault(uri string, options KeyVaultOptions) *Store {
	return &Store{
		uri: uri,
	}
}

func (s Store) Get(key string) (string, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", err
	}

	client, err := azsecrets.NewClient(s.uri, cred, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.GetSecret(context.Background(), key, "", nil)
	if err != nil {
		return "", err
	}

	return *resp.Value, nil
}
