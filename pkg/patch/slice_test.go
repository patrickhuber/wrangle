package patch_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/wrangle/pkg/patch"
	"github.com/stretchr/testify/require"
)

func TestSlice(t *testing.T) {
	t.Run("can add", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		update := patch.NewSlice(
			patch.SliceAppend(10))
		val, ok := update.Apply(reflect.ValueOf(slice))
		require.True(t, ok)
		require.False(t, val.IsNil())
		require.Equal(t, len(slice)+1, val.Len())
	})
	t.Run("can remove at", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		update := patch.NewSlice(
			patch.SliceRemoveAt(2))
		val, ok := update.Apply(reflect.ValueOf(slice))
		require.True(t, ok)
		require.False(t, val.IsNil())
		require.Equal(t, len(slice)-1, val.Len())
		require.Equal(t, 4, val.Index(2).Interface())

	})
	t.Run("can remove", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		length := len(slice)
		update := patch.NewSlice(
			patch.SliceRemove(func(v reflect.Value) bool {
				i := v.Int()
				return i == int64(8)
			}))
		val, ok := update.Apply(reflect.ValueOf(slice))
		require.True(t, ok)
		require.False(t, val.IsNil())
		require.Equal(t, length-1, val.Len())
	})
	t.Run("can modify", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5, 6, 7, 8}
		update := patch.NewSlice(
			patch.SliceModifyAt(5, 10))
		val, ok := update.Apply(reflect.ValueOf(slice))
		require.True(t, ok)
		require.False(t, val.IsNil())
		require.Equal(t, len(slice), val.Len())
		require.Equal(t, 10, val.Index(5).Interface())
	})
}
