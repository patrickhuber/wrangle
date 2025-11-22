package vault_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/vault"
	"github.com/stretchr/testify/require"
)

const (
	integrationEnv     = "INTEGRATION"
	vaultAddrEnv       = "INTEGRATION_VAULT_ADDR"
	vaultTokenEnv      = "INTEGRATION_VAULT_TOKEN"
	vaultRoleIDEnv     = "INTEGRATION_VAULT_ROLE_ID"
	vaultSecretIDEnv   = "INTEGRATION_VAULT_SECRET_ID"
)

// runStoreTests runs common store tests for the given store instance
func runStoreTests(t *testing.T, store stores.Store, keyPrefix string) {
	t.Run("set", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: keyPrefix + "-test",
		}}
		value := keyPrefix + "-value"
		err := store.Set(key, value)
		require.NoError(t, err)

		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, value, v)
	})

	t.Run("get", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: keyPrefix + "-test",
		}}
		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, keyPrefix+"-value", v)
	})

	t.Run("list", func(t *testing.T) {
		items, err := store.List()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(items), 1)
	})
}

func TestVaultWithToken(t *testing.T) {
	if os.Getenv(integrationEnv) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", integrationEnv)
	}

	address := os.Getenv(vaultAddrEnv)
	if address == "" {
		t.Skipf("skipping integration tests: set %s environment variable", vaultAddrEnv)
	}

	token := os.Getenv(vaultTokenEnv)
	if token == "" {
		t.Skipf("skipping integration tests: set %s environment variable", vaultTokenEnv)
	}

	factory := vault.NewFactory()
	store, err := factory.Create(map[string]string{
		"address": address,
		"token":   token,
		"path":    "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	runStoreTests(t, store, "token")
}

func TestVaultWithAppRole(t *testing.T) {
	if os.Getenv(integrationEnv) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", integrationEnv)
	}

	address := os.Getenv(vaultAddrEnv)
	if address == "" {
		t.Skipf("skipping integration tests: set %s environment variable", vaultAddrEnv)
	}

	roleID := os.Getenv(vaultRoleIDEnv)
	if roleID == "" {
		t.Skipf("skipping integration tests: set %s environment variable", vaultRoleIDEnv)
	}

	secretID := os.Getenv(vaultSecretIDEnv)
	if secretID == "" {
		t.Skipf("skipping integration tests: set %s environment variable", vaultSecretIDEnv)
	}

	factory := vault.NewFactory()
	store, err := factory.Create(map[string]string{
		"address":   address,
		"role_id":   roleID,
		"secret_id": secretID,
		"path":      "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	runStoreTests(t, store, "approle")
}
