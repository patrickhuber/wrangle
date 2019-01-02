package credhub

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/patrickhuber/wrangle/store"

	credhubcli "code.cloudfoundry.org/credhub-cli/credhub"
	"code.cloudfoundry.org/credhub-cli/credhub/auth"
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

func (s *credHubStore) Get(key string) (store.Data, error) {
	ch := s.credHub
	i := -1
	for i = len(key) - 1; i >= 0; i-- {
		if key[i] == '.' {
			break
		}
		if key[i] == '/' {
			i = -1
			break
		}
	}
	property := ""
	if i > 0 {
		property = key[i+1 : len(key)]
		key = key[0:i]
	}
	cred, err := ch.GetLatestVersion(key)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to lookup credential '%s'.", key)
	}

	if property == "" {
		return store.NewData(
			cred.Name,
			cred.Value), nil
	}

	m, ok := cred.Value.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unable to find property '%s' in credential '%s'", property, key)
	}
	propertyValue, ok := m[property]
	if !ok {
		return nil, fmt.Errorf("unable to find property '%s' in credential '%s'", property, key)
	}

	return store.NewData(
		cred.Name,
		propertyValue,
	), nil
}

func (s *credHubStore) Delete(key string) error {
	return fmt.Errorf("not implemented")
}

func (s *credHubStore) Type() string {
	return "credhub"
}

func (s *credHubStore) Set(key string, value string) (string, error) {
	ch := s.credHub
	_, err := ch.SetCredential(key, "value", value)
	if err != nil {
		return "", err
	}
	return value, nil
}

func (s *credHubStore) String() string {
	return s.Name()
}
