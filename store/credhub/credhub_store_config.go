package credhub

import (
	"strconv"

	"github.com/patrickhuber/wrangle/config"
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

func NewCredHubStoreConfig(configSource *config.Store) (*CredHubStoreConfig, error) {

	credHubStoreConfig := &CredHubStoreConfig{
		SkipTLSValidation: false,
	}

	credHubStoreConfig.Name = configSource.Name
	if value, ok := configSource.Params["skip_tls_validation"]; ok {
		boolValue, err := strconv.ParseBool(value)
		if err != nil {
			return nil, err
		}
		credHubStoreConfig.SkipTLSValidation = boolValue
	}

	if value, ok := configSource.Params["ca_cert"]; ok {
		credHubStoreConfig.CaCert = value
	}

	if value, ok := configSource.Params["username"]; ok {
		credHubStoreConfig.Username = value
	}

	if value, ok := configSource.Params["password"]; ok {
		credHubStoreConfig.Password = value
	}

	if value, ok := configSource.Params["client_id"]; ok {
		credHubStoreConfig.ClientID = value
	}

	if value, ok := configSource.Params["client_secret"]; ok {
		credHubStoreConfig.ClientSecret = value
	}

	if value, ok := configSource.Params["server"]; ok {
		credHubStoreConfig.Server = value
	}

	return credHubStoreConfig, nil
}
