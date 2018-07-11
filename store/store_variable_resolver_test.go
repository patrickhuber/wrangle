package store_test

import (
	"testing"

	"github.com/patrickhuber/wrangle/store"
	"github.com/patrickhuber/wrangle/store/memory"
	"github.com/stretchr/testify/require"
)

func TestVariableStoreResolver(t *testing.T) {

	t.Run("CanGetValueFromResolver", func(t *testing.T) {
		r := require.New(t)
		memoryStore := memory.NewMemoryStore("test")
		_, err := memoryStore.Put("key", "value")
		r.Nil(err)
		resolver := store.NewStoreVariableResolver(memoryStore)
		value, err := resolver.Get("key")
		r.Nil(err)
		r.Equal("value", value)
	})

}
