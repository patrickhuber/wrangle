package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

	if resp.ContentType != nil && *resp.ContentType == "application/json" {
		return s.decodeJSON(*resp.Value)
	}
	return *resp.Value, nil
}

func (s Store) decodeJSON(str string) (any, error) {
	dec := json.NewDecoder(strings.NewReader(str))
	t, err := dec.Token()
	if err != nil {
		return nil, err
	}
	if d, ok := t.(json.Delim); ok {
		switch d {
		case '[':
			var array []any
			err = json.Unmarshal([]byte(str), &array)
			return array, err
		case '{':
			var object map[string]any
			err = json.Unmarshal([]byte(str), &object)
			return object, err
		default:
			return nil, fmt.Errorf("unexpected json delimiter when decoding secret json %c", d)
		}
	}
	return nil, fmt.Errorf("unexpected input when decoding secret json. Input does not represent a JSON object or array")
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
