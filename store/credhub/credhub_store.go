package credhub

import (
	"errors"
	"fmt"

	"github.com/patrickhuber/cli-mgr/store"

	credhubcli "github.com/cloudfoundry-incubator/credhub-cli/credhub"
	"github.com/cloudfoundry-incubator/credhub-cli/credhub/auth"
)

type credHubStore struct {
	name    string
	credHub *credhubcli.CredHub
}

func NewCredHubStore(config *CredHubStoreConfig) (*credHubStore, error) {
	if config.ClientID == "" {
		return nil, errors.New("ClientID is required")
	}
	if config.ClientSecret == "" {
		return nil, errors.New("ClientSecret is required")
	}
	if config.Server == "" {
		return nil, errors.New("Server is required")
	}

	options := createOptions(config)
	ch, err := credhubcli.New(config.Server, options...)
	if err != nil {
		return nil, err
	}
	return &credHubStore{
		credHub: ch,
		name:    config.Name,
	}, nil
}

func createOptions(config *CredHubStoreConfig) []credhubcli.Option {
	options := []credhubcli.Option{}
	options = append(options, credhubcli.SkipTLSValidation(config.SkipTLSValidation))
	if config.CaCert != "" {
		options = append(options, credhubcli.CaCerts(config.CaCert))
	}
	options = append(options, credhubcli.Auth(
		auth.UaaClientCredentials(
			config.ClientID,
			config.ClientSecret)))
	return options
}

func (s *credHubStore) Name() string {
	return s.name
}

func (s *credHubStore) GetByName(name string) (store.Data, error) {
	ch := s.credHub
	i := -1
	for i = len(name) - 1; i >= 0; i-- {
		if name[i] == '.' {
			break
		}
		if name[i] == '/' {
			i = -1
			break
		}
	}
	property := ""
	if i > 0 {
		property = name[i+1 : len(name)]
		name = name[0:i]
	}
	cred, err := ch.GetLatestVersion(name)
	if err != nil {
		return nil, err
	}

	if property == "" {
		return store.NewData(
			cred.Id,
			cred.Name,
			cred.Value), nil
	}

	m, ok := cred.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to find property '%s' in credential '%s'", property, name)
	}
	propertyValue, ok := m[property]
	if !ok {
		return nil, fmt.Errorf("unable to find property '%s' in credential '%s'", property, name)
	}

	return store.NewData(
		cred.Id,
		cred.Name,
		propertyValue,
	), nil
}

func (s *credHubStore) Delete(name string) (int, error) {
	return 0, fmt.Errorf("not implemented")
}

func (s *credHubStore) Type() string {
	return "credhub"
}

func (s *credHubStore) Put(name string, value string) (string, error) {
	ch := s.credHub
	_, err := ch.SetCredential(name, "value", value, credhubcli.Overwrite)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *credHubStore) String() string {
	return s.Name()
}
