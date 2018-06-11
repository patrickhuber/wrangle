//+build integration

package credhub

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCredHubIntegration(t *testing.T) {
	t.Run("CanConnectToCredhub", func(t *testing.T) {
		r := require.New(t)
		server := getAndCheckEnvironmentVariable(r, "CREDHUB_SERVER")
		clientSecret := getAndCheckEnvironmentVariable(r, "CREDHUB_SECRET")
		clientID := getAndCheckEnvironmentVariable(r, "CREDHUB_CLIENT")
		ca := getAndCheckEnvironmentVariable(r, "CREDHUB_CA_CERT")

		// TODO: pull creds from environment
		storeConfig := &CredHubStoreConfig{
			Server:       server,
			ClientSecret: clientSecret,
			ClientID:     clientID,
			CaCert:       ca,
		}
		store, err := NewCredHubStore(storeConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		echo, err := store.Put("/test", "test")
		r.Nil(err)
		r.Equal("test", echo)
		value, err := store.GetByName("/test")
		r.Nil(err)
		r.Equal(echo, value.Value)
	})
}

func getAndCheckEnvironmentVariable(r *require.Assertions, key string) string {

	value, ok := os.LookupEnv(key)
	r.Truef(ok, "missing environment variable '%s'", key)
	return value
}
