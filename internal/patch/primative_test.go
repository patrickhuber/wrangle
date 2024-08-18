package patch_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/wrangle/internal/patch"
	"github.com/stretchr/testify/require"
)

func TestPrimative(t *testing.T) {
	t.Run("string update new value", func(t *testing.T) {
		update := patch.NewString("new")
		v, ok := update.Apply(reflect.ValueOf("old"))
		require.True(t, ok)
		require.Equal(t, "new", v.String())
	})
	t.Run("string update no change", func(t *testing.T) {
		update := patch.NewString("old")
		_, ok := update.Apply(reflect.ValueOf("old"))
		require.False(t, ok)
	})
	t.Run("int update change", func(t *testing.T) {
		update := patch.NewInt(10)
		v, ok := update.Apply(reflect.ValueOf(20))
		require.True(t, ok)
		require.EqualValues(t, 10, v.Int())
	})
	t.Run("int update no change", func(t *testing.T) {
		update := patch.NewInt(10)
		_, ok := update.Apply(reflect.ValueOf(10))
		require.False(t, ok)
	})
	t.Run("bool update sets true", func(t *testing.T) {
		update := patch.NewBool(false)
		v, ok := update.Apply(reflect.ValueOf(true))
		require.True(t, ok)
		require.False(t, v.Bool())
	})
	t.Run("bool update no change", func(t *testing.T) {
		update := patch.NewBool(true)
		_, ok := update.Apply(reflect.ValueOf(true))
		require.False(t, ok)
	})
}
