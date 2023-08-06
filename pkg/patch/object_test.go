package patch_test

import (
	"reflect"
	"testing"

	"github.com/patrickhuber/wrangle/pkg/patch"
	"github.com/stretchr/testify/require"
)

type Car struct {
	Make  string
	Model string
	Year  int
}

type Parent struct {
	Name     string
	Child    Child
	ChildPtr *Child
}

type Child struct {
	Name string
}

func TestObject(t *testing.T) {
	t.Run("can modify", func(t *testing.T) {
		car := &Car{
			Make:  "Ford",
			Model: "F150",
			Year:  2000,
		}
		update := patch.NewObject(
			patch.ObjectSetField("Make", "Tesla"))

		_, ok := update.Apply(reflect.ValueOf(car))
		require.True(t, ok)
		require.Equal(t, "Tesla", car.Make)
	})
	t.Run("can set child field", func(t *testing.T) {
		parent := &Parent{
			Name: "parent",
			Child: Child{
				Name: "child",
			},
		}
		update := patch.NewObject(
			patch.ObjectSetField("Child",
				patch.NewObject(
					patch.ObjectSetField("Name", "test"))))

		_, ok := update.Apply(reflect.ValueOf(parent))
		require.True(t, ok)
		require.Equal(t, "test", parent.Child.Name)
	})
	t.Run("can set child ptr", func(t *testing.T) {
		parent := &Parent{
			Name: "parent",
		}
		update := patch.NewObject(
			patch.ObjectSetField("ChildPtr",
				patch.NewObject(
					patch.ObjectSetField("Name", "test"))))

		_, ok := update.Apply(reflect.ValueOf(parent))
		require.True(t, ok)
		require.NotNil(t, parent.ChildPtr)
		require.Equal(t, "test", parent.ChildPtr.Name)
	})
}
