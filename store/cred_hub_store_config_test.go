package store

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/stretchr/testify/require"
)

func TestCredHubStoreConfig(t *testing.T) {
	t.Run("CanCreateCredHubStoreConfig", func(t *testing.T) {
		certificate := "-----BEGIN CERTIFICATE-----\n...\n-----END CERTIFICATE-----"

		name := "name"
		configSourceType := "type"

		// params
		server := "server"
		username := "username"
		password := "password"
		clientID := "client_id"
		clientSecret := "client_secret"

		configSource := &config.ConfigSource{
			Name:             name,
			ConfigSourceType: configSourceType,
			Params: map[string]string{
				"ca_cert":       certificate,
				"username":      username,
				"password":      password,
				"client_id":     clientID,
				"client_secret": clientSecret,
				"server":        server,
			},
		}

		storeConfig, err := NewCredHubStoreConfig(configSource)
		require := require.New(t)
		require.Nil(err)
		require.Equal(certificate, storeConfig.CaCert)
		require.Equal(username, storeConfig.Username)
		require.Equal(password, storeConfig.Password)
		require.Equal(clientID, storeConfig.ClientID)
		require.Equal(clientSecret, storeConfig.ClientSecret)
		require.Equal(server, storeConfig.Server)
		require.Equal(false, storeConfig.SkipTLSValidation)
	})
}
