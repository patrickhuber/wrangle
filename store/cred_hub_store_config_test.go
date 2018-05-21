package store

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
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
		if err != nil {
			t.Error(err)
			return
		}
		if certificate != storeConfig.CaCert {
			t.Errorf("invalid certificate. expected '%s' found '%s'", certificate, storeConfig.CaCert)
			return
		}
		if username != storeConfig.Username {
			t.Errorf("invalid username. expected '%s' found '%s'", username, storeConfig.Username)
			return
		}
		if password != storeConfig.Password {
			t.Errorf("invalid password. expected '%s' found '%s'", password, storeConfig.Password)
			return
		}
		if clientID != storeConfig.ClientID {
			t.Errorf("invalid client_id. expected '%s' found '%s'", clientID, storeConfig.ClientID)
			return
		}
		if clientSecret != storeConfig.ClientSecret {
			t.Errorf("invalid client_secret. expected '%s' found '%s'", clientSecret, storeConfig.ClientSecret)
			return
		}
		if server != storeConfig.Server {
			t.Errorf("invalid server. expected '%s' found '%s'", server, storeConfig.Server)
			return
		}
		if false != storeConfig.SkipTLSValidation {
			t.Errorf("expected SkipTLSValidation flag of '%v' found '%v'", false, storeConfig.SkipTLSValidation)
			return
		}
	})
}
