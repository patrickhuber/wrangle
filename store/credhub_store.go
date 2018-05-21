package store

import (
	"fmt"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

type CredHubStore struct {
	Name    string
	CredHub *credhub.CredHub
}

func NewCredHubStore(config *CredHubStoreConfig) (*CredHubStore, error) {
	options := createOptions(config)
	ch, err := credhub.New(config.Server, options...)
	if err != nil {
		return nil, err
	}
	return &CredHubStore{
		CredHub: ch,
		Name:    config.Name,
	}, nil
}

func createOptions(config *CredHubStoreConfig) []credhub.Option {
	options := []credhub.Option{}
	options = append(options, credhub.SkipTLSValidation(config.SkipTLSValidation))
	options = append(options, credhub.Auth(
		auth.UaaClientCredentials(
			config.ClientID,
			config.ClientSecret)))
	return []credhub.Option{}
}

func (store *CredHubStore) GetName() string {
	return store.Name
}

func (store *CredHubStore) GetByName(name string) (StoreData, error) {
	ch := store.CredHub
	cred, err := ch.GetLatestVersion(name)
	if err != nil {
		return StoreData{}, err
	}
	return StoreData{
		ID:    cred.Id,
		Name:  cred.Name,
		Value: cred.Value,
	}, nil
}

func (store *CredHubStore) Delete(name string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (store *CredHubStore) GetType() string {
	return "credhub"
}

func (store *CredHubStore) Put(name string, value string) (string, error) {
	return value, fmt.Errorf("not implemented")
}
