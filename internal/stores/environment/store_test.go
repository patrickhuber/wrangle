package environment_test

import (
	"testing"

	"github.com/patrickhuber/go-cross"
	"github.com/patrickhuber/go-cross/arch"
	"github.com/patrickhuber/go-cross/platform"
	"github.com/patrickhuber/wrangle/internal/stores"
	"github.com/patrickhuber/wrangle/internal/stores/environment"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentStore(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	env := target.Env()
	require.NoError(t, env.Set("MY_VAR", "value"))

	store, err := environment.NewFactory(env).Create(nil)
	require.NoError(t, err)
	key, err := stores.ParseKey("MY_VAR")
	require.NoError(t, err)

	value, found, err := store.Get(key)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "value", value)

	require.ErrorContains(t, store.Set(key, "updated"), "read-only")
	require.Equal(t, "value", env.Get("MY_VAR"))

	keys, err := store.List()
	require.NoError(t, err)
	require.Contains(t, keys, stores.Key{
		Data: stores.Data{Name: "MY_VAR", Version: stores.Version{Latest: true}},
	})
}

func TestEnvironmentStoreMissingVariableAndInvalidKeys(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	store, err := environment.NewFactory(target.Env()).Create(nil)
	require.NoError(t, err)

	missing, err := stores.ParseKey("DOES_NOT_EXIST")
	require.NoError(t, err)
	_, found, err := store.Get(missing)
	require.NoError(t, err)
	require.False(t, found)

	for _, invalid := range []string{"", "env/one", "MY_VAR@1"} {
		t.Run(invalid, func(t *testing.T) {
			key, err := stores.ParseKey(invalid)
			require.NoError(t, err)
			_, found, err := store.Get(key)
			require.NoError(t, err)
			require.False(t, found)
		})
	}
}

func TestEnvironmentStoreUnspecifiedVersion(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	require.NoError(t, target.Env().Set("MY_VAR", "value"))
	store, err := environment.NewFactory(target.Env()).Create(nil)
	require.NoError(t, err)

	value, found, err := store.Get(stores.Key{Data: stores.Data{Name: "MY_VAR"}})
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, "value", value)
}

func TestEnvironmentStoreIncludeAndExcludePatterns(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	env := target.Env()
	require.NoError(t, env.Set("APP_TOKEN", "app-token"))
	require.NoError(t, env.Set("BUILD_NUMBER", "123"))
	require.NoError(t, env.Set("SECRET_TOKEN", "secret"))
	require.NoError(t, env.Set("UNRELATED", "value"))

	store, err := environment.NewFactory(env).Create(map[string]any{
		"include": []any{"*_TOKEN", "BUILD_*"},
		"exclude": []any{"SECRET_*"},
	})
	require.NoError(t, err)

	for _, name := range []string{"APP_TOKEN", "BUILD_NUMBER"} {
		key, err := stores.ParseKey(name)
		require.NoError(t, err)
		value, found, err := store.Get(key)
		require.NoError(t, err)
		require.True(t, found)
		require.NotEmpty(t, value)
	}

	for _, name := range []string{"SECRET_TOKEN", "UNRELATED"} {
		key, err := stores.ParseKey(name)
		require.NoError(t, err)
		_, found, err := store.Get(key)
		require.NoError(t, err)
		require.False(t, found)
	}

	keys, err := store.List()
	require.NoError(t, err)
	require.Len(t, keys, 2)
	require.Contains(t, keys, stores.Key{
		Data: stores.Data{Name: "APP_TOKEN", Version: stores.Version{Latest: true}},
	})
	require.Contains(t, keys, stores.Key{
		Data: stores.Data{Name: "BUILD_NUMBER", Version: stores.Version{Latest: true}},
	})
}

func TestEnvironmentStoreEmptyIncludeExposesNoVariables(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	require.NoError(t, target.Env().Set("MY_VAR", "value"))

	store, err := environment.NewFactory(target.Env()).Create(map[string]any{
		"include": []any{},
	})
	require.NoError(t, err)

	key, err := stores.ParseKey("MY_VAR")
	require.NoError(t, err)
	_, found, err := store.Get(key)
	require.NoError(t, err)
	require.False(t, found)

	keys, err := store.List()
	require.NoError(t, err)
	require.Empty(t, keys)
}

func TestEnvironmentStoreExcludeOnlyFiltersDefaultAll(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)
	env := target.Env()
	require.NoError(t, env.Set("KEEP", "value"))
	require.NoError(t, env.Set("SKIP_SECRET", "secret"))

	store, err := environment.NewFactory(env).Create(map[string]any{
		"exclude": []any{"*_SECRET"},
	})
	require.NoError(t, err)

	keep, err := stores.ParseKey("KEEP")
	require.NoError(t, err)
	_, found, err := store.Get(keep)
	require.NoError(t, err)
	require.True(t, found)

	skip, err := stores.ParseKey("SKIP_SECRET")
	require.NoError(t, err)
	_, found, err = store.Get(skip)
	require.NoError(t, err)
	require.False(t, found)
}

func TestEnvironmentStoreRejectsInvalidPatterns(t *testing.T) {
	target := cross.NewTest(platform.Linux, arch.AMD64)

	_, err := environment.NewFactory(target.Env()).Create(map[string]any{
		"include": []any{"["},
	})
	require.ErrorContains(t, err, "invalid glob")
}
