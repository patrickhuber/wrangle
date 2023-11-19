package stores_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/stores"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	t.Run("store", func(t *testing.T) {
		str := "store:"
		key, err := stores.Parse(str)
		require.NoError(t, err)
		require.Equal(t, "store", key.Store)
	})
	t.Run("store_secret", func(t *testing.T) {
		str := "store:secret"
		key, err := stores.Parse(str)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.Equal(t, "store", key.Store)
		require.NotNil(t, key.Secret, "secret is nil")
		require.Equal(t, "secret", key.Secret.Name)
	})
	t.Run("store_secret_version", func(t *testing.T) {
		str := "store:secret@1"
		key, err := stores.Parse(str)
		require.NoError(t, err)
		require.NotNil(t, key)
		require.Equal(t, "store", key.Store)
		require.NotNil(t, key.Secret, "secret is nil")
		require.Equal(t, "secret", key.Secret.Name)
		require.Equal(t, "1", key.Secret.Version)
	})
	t.Run("long_path", func(t *testing.T) {
		// {store}:{secret}/{path_item}/{path_item}/{path_item}/{path_item}
		str := "store:secret/this/is/the/path"
		key, err := stores.Parse(str)
		require.NoError(t, err)
		require.Equal(t, "store", key.Store)
		require.Equal(t, []stores.PathItem{
			stores.Name{Value: "this"},
			stores.Name{Value: "is"},
			stores.Name{Value: "the"},
			stores.Name{Value: "path"},
		}, key.Path)
	})
}
