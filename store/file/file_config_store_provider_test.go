package file

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/stretchr/testify/require"
)

func TestFileConfigStoreProvider(t *testing.T) {
	t.Run("CanGetByName", func(t *testing.T) {
		r := require.New(t)
		provider := FileConfigStoreProvider{}
		name := provider.GetName()
		r.Equal("file", name)
	})

	t.Run("CanCreate", func(t *testing.T) {
		r := require.New(t)
		provider := FileConfigStoreProvider{}
		store, err := provider.Create(&config.ConfigSource{Name: ""})
		r.Nil(err)
		r.NotNil(store)
	})
}
