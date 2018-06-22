package store_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	patch "github.com/cppforlife/go-patch/patch"
)

func TestGoPatch(t *testing.T) {
	t.Run("CanFindKeyValue", func(t *testing.T) {
		require := require.New(t)
		pointer, err := patch.NewPointerFromString("/key1")
		if err != nil {
			t.Error(err)
			return
		}
		doc := map[interface{}]interface{}{
			"key1": "abc",
			"key2": "xyz",
		}
		response, err := patch.FindOp{Path: pointer}.Apply(doc)
		require.Nil(err)
		require.Equal("abc", response)
	})

	t.Run("CanCreatePointer", func(t *testing.T) {
		require := require.New(t)
		ptr, err := patch.NewPointerFromString("/some/path")
		require.Nil(err)
		require.Equal(3, len(ptr.Tokens()))
	})
}
