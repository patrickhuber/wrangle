package patch_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/wrangle/internal/patch"
	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	setup := func(t *testing.T) map[string]int {
		return map[string]int{
			"one":   1,
			"two":   2,
			"three": 3,
		}
	}
	t.Run("can add", func(t *testing.T) {
		m := setup(t)
		p := patch.NewMap(
			patch.MapSet("four", 4))

		result, ok := p.Apply(reflect.ValueOf(m))
		require.True(t, ok)
		require.False(t, result.IsNil())
		require.Equal(t, 4, len(result.MapKeys()))

		v, ok := m["four"]
		require.True(t, ok)
		require.Equal(t, 4, v)
	})
	t.Run("can remove", func(t *testing.T) {
		m := setup(t)
		p := patch.NewMap(
			patch.MapRemove("two"))

		result, ok := p.Apply(reflect.ValueOf(m))
		require.True(t, ok)
		require.False(t, result.IsNil())
		require.Equal(t, 2, len(result.MapKeys()))

		_, ok = m["two"]
		require.False(t, ok)
	})
	t.Run("can set", func(t *testing.T) {
		m := setup(t)
		p := patch.NewMap(
			patch.MapSet("two", 4))

		result, ok := p.Apply(reflect.ValueOf(m))
		require.True(t, ok)
		require.False(t, result.IsNil())
		require.Equal(t, 3, len(result.MapKeys()))

		v, ok := m["two"]
		require.True(t, ok)
		require.Equal(t, 4, v)
	})
}
