package memory

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/config"
	"github.com/stretchr/testify/require"
)

func TestMemoryStoreProvider(t *testing.T) {
	t.Run("CanCreateMemoryStore", func(t *testing.T) {
		r := require.New(t)
		provider := NewMemoryStoreProvider()
		name := provider.GetName()
		r.Equal("memory", name)
		store, err := provider.Create(&config.ConfigSource{})
		r.Nil(err)
		r.NotNil(store)
	})
}
