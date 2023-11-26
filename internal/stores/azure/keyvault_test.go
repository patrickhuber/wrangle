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
		d, err := s.Get(stores.Key{Data: &stores.Data{Name: "test"}})
		require.NoError(t, err)
		str, ok := d.(string)
		require.True(t, ok)
		require.Equal(t, str, "value")
	})

	t.Run("lookup", func(t *testing.T) {
		s := azure.NewKeyVault(uri, nil)
		d, ok, err := s.Lookup(stores.Key{Data: &stores.Data{Name: "test"}})
		require.NoError(t, err)
		require.True(t, ok)
		str, ok := d.(string)
		require.True(t, ok)
		require.Equal(t, str, "value")
	})
}
