package keyring_test

import (
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/keyring"
	"github.com/stretchr/testify/require"
)

func TestKeyRing(t *testing.T) {
	var INTEGRATION = "INTEGRATION"
	if os.Getenv(INTEGRATION) == "" {
		t.Skipf("skipping integration tests: set %s environment variable", INTEGRATION)
	}
	factory := keyring.NewFactory()
	vault, err := factory.Create(map[string]string{"service": "test"})
	require.Nil(t, err)
	require.NotNil(t, vault)

	t.Run("set", func(t *testing.T) {
		key := stores.Key{Data: stores.Data{
			Name: "test",
		}}
		err := vault.Set(key, "test")
		require.Nil(t, err)

		v, ok, err := vault.Get(key)
		require.Nil(t, err)
		require.True(t, ok)
		require.NotNil(t, v)
		require.Equal(t, "test", v)

		items, err := vault.List()
		require.NoError(t, err)
		require.Equal(t, 1, len(items))
	})
}
