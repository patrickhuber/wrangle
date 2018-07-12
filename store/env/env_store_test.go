package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEnvStore(t *testing.T) {

	t.Run("CanReadEnvironmentVariable", func(t *testing.T) {
		r := require.New(t)
		err := os.Setenv("TEST123", "abc123")
		r.Nil(err)
		lookup := map[string]string{
			"somevalue": "TEST123",
		}

		store := NewEnvStore("", lookup)
		r.NotNil(store)

		data, err := store.GetByName("somevalue")
		r.Nil(err)
		r.NotNil(data)
		r.Equal("abc123", data.Value())
	})

	t.Run("CanReadEnvironmentVariableWithPrefixedName", func(t *testing.T) {
		r := require.New(t)
		err := os.Setenv("TEST123", "abc123")

		r.Nil(err)
		lookup := map[string]string{
			"somevalue": "TEST123",
		}

		store := NewEnvStore("", lookup)
		r.NotNil(store)

		data, err := store.GetByName("/somevalue")
		r.Nil(err)
		r.NotNil(data)
		r.Equal("abc123", data.Value())
	})

	t.Run("ErrorsIfEnvironmentVariableNotSet", func(t *testing.T) {
		r := require.New(t)
		lookup := map[string]string{
			"somevalue": "TEST123",
		}
		store := NewEnvStore("", lookup)
		r.NotNil(store)

		os.Unsetenv("TEST123")
		_, err := store.GetByName("somevalue")
		r.NotNil(err)
	})

	t.Run("CanGetStoreName", func(t *testing.T) {
		r := require.New(t)
		store := NewEnvStore("env", nil)
		r.NotNil(store)
		r.Equal("env", store.Name())
	})

	t.Run("CanGetStoreType", func(t *testing.T) {
		r := require.New(t)
		store := NewEnvStore("", nil)
		r.NotNil(store)
		r.Equal("env", store.Type())
	})
}
