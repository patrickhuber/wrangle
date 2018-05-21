package store

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/patrickhuber/cli-mgr/config"
)

func TestStoreFactory(t *testing.T) {
	storeFactory := NewStoreFactory()
	t.Run("CanCreateMemoryStore", func(t *testing.T) {
		require := require.New(t)

		configSource := &config.ConfigSource{
			ConfigSourceType: "memory",
			Config:           "config",
			Name:             "name",
			Params:           map[string]string{}}

		store, err := storeFactory.Create(configSource)
		require.Nil(err)
		require.NotNil(store)

		name := store.GetName()
		require.Equal(name, configSource.Name)
	})

	t.Run("CanCreateFileStore", func(t *testing.T) {
		require := require.New(t)

		configSource := &config.ConfigSource{
			ConfigSourceType: "file",
			Config:           "config",
			Name:             "name",
			Params:           map[string]string{}}

		store, err := storeFactory.Create(configSource)
		require.Nil(err)
		require.NotNil(store)

		name := store.GetName()
		require.Equal(name, configSource.Name)
	})
}
