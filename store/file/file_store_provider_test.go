package file

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/patrickhuber/wrangle/config"
	"github.com/stretchr/testify/require"
)

func TestFileStoreProvider(t *testing.T) {
	t.Run("CanGetByName", func(t *testing.T) {
		r := require.New(t)
		provider := NewFileStoreProvider(afero.NewMemMapFs())
		name := provider.GetName()
		r.Equal("file", name)
	})

	t.Run("CanCreate", func(t *testing.T) {
		r := require.New(t)
		provider := NewFileStoreProvider(afero.NewMemMapFs())
		configSource := &config.ConfigSource{
			Name:             "test",
			ConfigSourceType: "file",
			Params: map[string]string{
				"path": "/file",
			},
		}
		store, err := provider.Create(configSource)
		r.Nil(err)
		r.NotNil(store)
	})
}
