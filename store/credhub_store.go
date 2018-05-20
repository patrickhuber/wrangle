package store

import (
	"fmt"
	"strconv"

	"github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/patrickhuber/cli-mgr/config"
)

type CredHubStoreConfig struct {
	Name              string
	Username          string
	Password          string
	Server            string
	ClientSecret      string
	ClientID          string
	CaCert            string
	SkipTLSValidation bool
}

func NewCredHubStoreConfig(configSource *config.ConfigSource) (*CredHubStoreConfig, error) {
	skipTLSValidation, err := strconv.ParseBool(configSource.Params["skipTLSValidation"])
	if err != nil {
		return nil, err
	}
	credHubStoreConfig := &CredHubStoreConfig{
		CaCert:            configSource.Params["caCert"],
		ClientID:          configSource.Params["clientID"],
		ClientSecret:      configSource.Params["clientSecret"],
		Username:          configSource.Params["userName"],
		Password:          configSource.Params["password"],
		SkipTLSValidation: skipTLSValidation,
	}
	return credHubStoreConfig, nil
}

type CredHubStore struct {
	Name    string
	CredHub *credhub.CredHub
}

func NewCredHubStore(config *CredHubStoreConfig) (*CredHubStore, error) {
	options := []credhub.Option{}
	// create options here
	ch, err := credhub.New(config.Server, options...)
	if err != nil {
		return nil, err
	}
	return &CredHubStore{
		CredHub: ch,
		Name:    config.Name,
	}, nil
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
	switch cred.Type {
	case "":
		break
	}
	return StoreData{}, nil
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
