package templates

import (
	"testing"

	"github.com/patrickhuber/cli-mgr/store/memory"
	"github.com/stretchr/testify/require"
)

func TestChainReslover(t *testing.T) {
	t.Run("CanReadFromSecondaryResolver", func(t *testing.T) {
		r := require.New(t)
		delegateResolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		r.NotNil(delegateResolver)
		store := memory.NewMemoryStore("")
		store.Put("key1", "((key))")
		resolver := NewChainVariableResolver(store, delegateResolver)
		data, err := resolver.Get("key1")
		r.Nil(err)
		r.Equal("value", data)
	})
	t.Run("CanChainMultipleResolvers", func(t *testing.T) {
		r := require.New(t)
		firstStore := memory.NewMemoryStore("first")
		firstStore.Put("one", "((two))")
		secondStore := memory.NewMemoryStore("second")
		secondStore.Put("two", "value")
		storeVariableResolver := NewStoreVariableResolver(secondStore)
		resolver := NewChainVariableResolver(firstStore, storeVariableResolver)
		data, err := resolver.Get("one")
		r.Nil(err)
		r.Equal("value", data)
	})
}
