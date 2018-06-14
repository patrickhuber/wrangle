package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChainReslover(t *testing.T) {
	t.Run("CanReadFromSecondaryResolver", func(t *testing.T) {
		r := require.New(t)
		delegateResolver, err := newSimpleResolver("key", "value")
		r.Nil(err)
		r.NotNil(delegateResolver)
		mainResolver, err := newSimpleResolver("key1", "((key))")
		r.Nil(err)
		r.NotNil(delegateResolver)
		resolver := NewChainVariableResolver(mainResolver, delegateResolver)
		data, err := resolver.Get("key1")
		r.Nil(err)
		r.Equal("value", data)
	})
	t.Run("CanChainMultipleResolvers", func(t *testing.T) {
		r := require.New(t)
		first, err := newSimpleResolver("one", "((two))")
		r.Nil(err)
		second, err := newSimpleResolver("two", "((three))")
		r.Nil(err)
		three, err := newSimpleResolver("three", "value")
		r.Nil(err)
		resolver := NewChainVariableResolver(first, second)
		resolver = NewChainVariableResolver(resolver, three)
		value, err := resolver.Get("one")
		r.Nil(err)
		r.Equal("value", value)
	})
}
