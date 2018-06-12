package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleResolver(t *testing.T) {
	t.Run("CanCreateSimpleResolver", func(t *testing.T) {
		r := require.New(t)
		resolver, err := newSimpleResolver("key", "value", "key1", "value1")
		r.Nil(err)
		r.NotNil(resolver)
		value, err := resolver.Get("key")
		r.Nil(err)
		r.NotNil(value)
		r.Equal("value", value)
		value, err = resolver.Get("key1")
		r.Nil(err)
		r.NotNil(value)
		r.Equal("value1", value)
	})
}
