package file

import (
	"testing"

	"github.com/spf13/afero"

	"github.com/patrickhuber/cli-mgr/config"
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
		store, err := provider.Create(&config.ConfigSource{Name: ""})
		r.Nil(err)
		r.NotNil(store)
	})
}
