package vault_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/vault"
	"github.com/stretchr/testify/require"
)

func TestVaultWithToken(t *testing.T) {
	var INTEGRATION = "INTEGRATION"
	var VAULT_ADDR_ENV = "INTEGRATION_VAULT_ADDR"
	var VAULT_TOKEN_ENV = "INTEGRATION_VAULT_TOKEN"

	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
	}

	address := os.Getenv(VAULT_ADDR_ENV)
	if address == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_ADDR_ENV)
	}

	token := os.Getenv(VAULT_TOKEN_ENV)
	if token == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_TOKEN_ENV)
	}

	factory := vault.NewFactory()
	store, err := factory.Create(map[string]string{
		"address": address,
		"token":   token,
		"path":    "secret",
	})
	require.NoError(t, err)
	require.NotNil(t, store)

	t.Run("set", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: "test",
		}}
		err := store.Set(key, "test-value")
		require.NoError(t, err)

		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, "test-value", v)
	})

	t.Run("get", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: "test",
		}}
		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, "test-value", v)
	})

	t.Run("list", func(t *testing.T) {
		items, err := store.List()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(items), 1)
	})
}

func TestVaultWithAppRole(t *testing.T) {
	var INTEGRATION = "INTEGRATION"
	var VAULT_ADDR_ENV = "INTEGRATION_VAULT_ADDR"
	var VAULT_ROLE_ID_ENV = "INTEGRATION_VAULT_ROLE_ID"
	var VAULT_SECRET_ID_ENV = "INTEGRATION_VAULT_SECRET_ID"

	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
	}

	address := os.Getenv(VAULT_ADDR_ENV)
	if address == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_ADDR_ENV)
	}

	roleID := os.Getenv(VAULT_ROLE_ID_ENV)
	if roleID == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_ROLE_ID_ENV)
	}

	secretID := os.Getenv(VAULT_SECRET_ID_ENV)
	if secretID == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_SECRET_ID_ENV)
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

	t.Run("set", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: "test-approle",
		}}
		err := store.Set(key, "approle-value")
		require.NoError(t, err)

		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, "approle-value", v)
	})

	t.Run("get", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: "test-approle",
		}}
		v, ok, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, "approle-value", v)
	})

	t.Run("list", func(t *testing.T) {
		items, err := store.List()
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(items), 1)
	})
}
