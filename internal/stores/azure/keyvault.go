package azure

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
	"github.com/patrickhuber/wrangle/internal/dataptr"
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

func (s Store) Get(key stores.Key) (any, bool, error) {
	client, err := s.client()
	if err != nil {
		return nil, false, err
	}

	version := ""
	if !key.Data.Version.Latest {
		version = key.Data.Version.Value
	}

	resp, err := client.GetSecret(context.Background(), key.Data.Name, version, nil)
	if err != nil {
		// check not found error
		return "", false, err
	}

	if resp.ContentType != nil && *resp.ContentType == "application/json" {
		v, err := s.decodeJSON(*resp.Value)
		return v, err == nil, err
	}
	return *resp.Value, true, nil
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

func (s Store) List() ([]stores.Key, error) {
	client, err := s.client()
	if err != nil {
		return nil, err
	}
	var keys []stores.Key

	resp := client.NewListSecretsPager(nil)
	latest := map[string]string{}
	for resp.More() {
		page, err := resp.NextPage(context.Background())
		if err != nil {
			return nil, err
		}
		for _, secret := range page.SecretListResult.Value {
			name := secret.ID.Name()
			version := secret.ID.Version()

			// cache the latest version
			if _, ok := latest[name]; !ok {
				resp, err := client.GetSecret(
					context.Background(),
					name,
					version, nil)
				if err != nil {
					return nil, err
				}
				latest[name] = resp.ID.Version()
			}

			keys = append(keys, stores.Key{
				Data: &stores.Data{
					Name: name,
					Version: stores.Version{
						Value: version,
					},
				},
				Path: &dataptr.DataPointer{},
			})
		}
	}
	return keys, nil
}

func (s Store) client() (*azsecrets.Client, error) {
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return nil, err
	}

	return azsecrets.NewClient(s.uri, cred, nil)
}
