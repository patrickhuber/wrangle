package azure_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/azure"
	"github.com/stretchr/testify/require"
)

func TestKeyVault(t *testing.T) {

	var uri string = "https://%s.vault.azure.net/"
	var INTEGRATION = "INTEGRATION"
	var VAULT_NAME_ENV = "INTEGRATION_AZURE_KEYVAULT_NAME"

	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
	}

	name := os.Getenv(VAULT_NAME_ENV)
	if name == "" {
		t.Skipf("skipping integration tests: set %s environment variable", VAULT_NAME_ENV)
	}

	uri = fmt.Sprintf(uri, name)

	t.Run("get", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		d, ok, err := s.Get(stores.Key{Data: stores.Data{Name: "test"}})
		require.NoError(t, err)
		require.True(t, ok)
		str, ok := d.(string)
		require.True(t, ok)
		require.Equal(t, str, "value")
	})

	t.Run("list", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		keys, err := s.List()
		require.NoError(t, err)
		require.Equal(t, 3, len(keys))
	})

	t.Run("get json object", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		d, ok, err := s.Get(stores.Key{Data: stores.Data{Name: "json-object"}})
		require.NoError(t, err)
		require.True(t, ok)
		// {"test":"value"}
		expected := map[string]any{
			"test": "value",
		}
		require.Equal(t, expected, d)
	})

	t.Run("get json array", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		d, ok, err := s.Get(stores.Key{Data: stores.Data{Name: "json-array"}})
		require.NoError(t, err)
		require.True(t, ok)
		// ["one", "two", "three"]
		expected := []any{"one", "two", "three"}
		require.Equal(t, expected, d)
	})

	t.Run("set", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		err := s.Set(stores.Key{Data: stores.Data{Name: "set-test"}}, "value")
		require.NoError(t, err)
	})
}
