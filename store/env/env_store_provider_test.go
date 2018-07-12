package env

import (
	"os"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/stretchr/testify/require"
)

func TestEnvStoreProvider(t *testing.T) {

	t.Run("CanCreateEnvStoreProvider", func(t *testing.T) {
		r := require.New(t)
		provider := NewEnvStoreProvider()
		name := provider.GetName()
		r.Equal("env", name)
	})

	t.Run("CanCreateEnvStore", func(t *testing.T) {
		r := require.New(t)
		provider := NewEnvStoreProvider()

		prop1 := "test1"
		prop2 := "test2"
		env1 := "TEST_WRANGLE_ENV_VAL1"
		env2 := "TEST_WRANGLE_ENV_VAL2"
		val1 := "abc123"
		val2 := "zyx987"

		source := &config.Store{
			Name:      "env",
			StoreType: "env",
			Params: map[string]string{
				prop1: env1,
				prop2: env2,
			},
		}

		// create the store
		s, err := provider.Create(source)
		r.Nil(err)
		r.NotNil(s)

		// make sure the name is set
		r.Equal(source.Name, s.Name())

		// set the env vars
		os.Setenv(env1, val1)
		os.Setenv(env2, val2)

		// verify
		d1, err := s.GetByName(prop1)
		r.Nil(err)
		r.Equal(val1, d1.Value())

		d2, err := s.GetByName(prop2)
		r.Nil(err)
		r.Equal(val2, d2.Value())
	})
}
